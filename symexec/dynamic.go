package symexec

import (
	"fmt"
	"go/types"
	"os"
	"runtime/debug"
	"strings"

	"github.com/aclements/go-z3/z3"
	"golang.org/x/tools/go/ssa"
)

const maxDepth = 100

func Dynamic() {
	os.Chdir("testdata")

	testcases, err := os.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, tc := range testcases {
		if tc.IsDir() || strings.HasSuffix(tc.Name(), "_test.go") {
			continue
		}
		r := AnalyzeFileDynamic(tc.Name())
		GenerateTests(tc.Name(), r)
	}
}

func AnalyzeFileDynamic(filename string) map[*ssa.Function][]Testcase {
	main := buildPackage(filename)
	res := make(map[*ssa.Function][]Testcase, 0)
	for _, v := range main.Members {
		if fn, ok := v.(*ssa.Function); ok && fn.Name() != "init" {
			res[fn] = dynamicFunction(fn, main)
		}
		if obj, ok := v.(*ssa.Type); ok {
			named := obj.Type().(*types.Named)
			n := named.NumMethods()
			for i := 0; i < n; i++ {
				fn := main.Prog.FuncValue(named.Method(i))
				res[fn] = dynamicFunction(fn, main)
			}
		}
	}
	return res
}

func dynamicFunction(fn *ssa.Function, pkg *ssa.Package) []Testcase {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[ERROR]", r)
			fmt.Println(string(debug.Stack()))
		}
	}()
	fmt.Println("::", "analyzing function", "'"+fn.Name()+"'")
	fmt.Println("::", "printing SSA blocks")
	printBlocks(fn)
	fmt.Println("::", "execute")
	return execute(fn, pkg, &RandomQueue{})
}

type State struct {
	nextFrameId int
	depth       int
	frames      []*Frame
}

func (s *State) copy() *State {
	stateCopy := &State{}
	for _, frame := range s.frames {
		stateCopy.frames = append(stateCopy.frames, frame.copy())
	}
	// patch in-progress dynamic calls
	for i := 0; i+1 < len(stateCopy.frames); i++ {
		caller := stateCopy.frames[i]
		callee := stateCopy.frames[i+1]
		if _, ok := caller.call.Body[len(caller.call.Body)-1].(*DynamicCall); !ok {
			panic("not a dynamic call")
		}
		caller.call.Body[len(caller.call.Body)-1] = callee.call
	}
	stateCopy.nextFrameId = s.nextFrameId
	stateCopy.depth = s.depth
	return stateCopy
}

func (s *State) currentFrame() *Frame {
	return s.frames[len(s.frames)-1]
}

func (s *State) formula() Formula {
	return And{s.frames[0].call.Body}
}

type Frame struct {
	id         int
	function   *ssa.Function
	blockOrder []int
	call       *DynamicCall
	nextBlock  int
	nextInstr  int
}

func (frame *Frame) push(f Formula) {
	frame.call.Body = append(frame.call.Body, f)
}

func (frame *Frame) newVar(reg Register) Var {
	tmp := &TempRegister{
		name: reg.Name(),
		t:    reg.Type(),
	}
	if frame.id > 0 {
		tmp.name = fmt.Sprintf("%d#%s", frame.id, tmp.name)
	}
	return NewVar(tmp)
}

func (frame *Frame) copy() *Frame {
	var blockOrder []int
	var body []Formula
	return &Frame{
		id:         frame.id,
		function:   frame.function,
		blockOrder: append(blockOrder, frame.blockOrder...),
		call: &DynamicCall{
			Result: frame.call.Result,
			Name:   frame.call.Name,
			Args:   frame.call.Args,
			Params: frame.call.Params,
			Body:   append(body, frame.call.Body...),
		},
		nextBlock: frame.nextBlock,
		nextInstr: frame.nextInstr,
	}
}

func execute(fn *ssa.Function, pkg *ssa.Package, queue Queue) []Testcase {
	var testcases []Testcase
	entryPoint := &DynamicCall{
		Result: Var{},
		Name:   fn.Name(),
		Args:   nil,
		Params: nil,
		Body:   nil,
	}
	entryFrame := &Frame{function: fn, call: entryPoint}
	queue.push(&State{frames: []*Frame{entryFrame}})
	for !queue.empty() {
		state := queue.pop()
		state.depth += 1
		if state.depth >= maxDepth {
			fmt.Println("[WARNING] max depth reached")
			continue
		}
		frame := state.currentFrame()
		block := frame.function.Blocks[frame.nextBlock]
	instructionLoop:
		for index, instr := range block.Instrs {
			if index < frame.nextInstr {
				continue
			}
			if index == 0 {
				frame.blockOrder = append(frame.blockOrder, frame.nextBlock)
			} else {
				frame.nextInstr = 0
			}
			switch v := instr.(type) {
			case *ssa.BinOp:
				frame.push(BinOp{
					Result: frame.newVar(v),
					Left:   frame.newVar(v.X),
					Op:     v.Op.String(),
					Right:  frame.newVar(v.Y),
				})
			case *ssa.If:
				{
					thenState := state.copy()
					thenFrame := thenState.currentFrame()
					thenFrame.nextBlock = v.Block().Succs[0].Index
					thenFrame.push(Condition{
						Cond:   frame.newVar(v.Cond),
						IsTrue: true,
					})
					if _, sat := solve(fn, thenState.formula()); sat {
						queue.push(thenState)
					}
				}
				{
					elseState := state.copy()
					elseFrame := elseState.currentFrame()
					elseFrame.nextBlock = v.Block().Succs[1].Index
					elseFrame.push(Condition{
						Cond:   frame.newVar(v.Cond),
						IsTrue: false,
					})
					if _, sat := solve(fn, elseState.formula()); sat {
						queue.push(elseState)
					}
				}
				break instructionLoop
			case *ssa.Jump:
				state.currentFrame().nextBlock = v.Block().Succs[0].Index
				queue.push(state)
				break instructionLoop
			case *ssa.Return:
				var results []Var
				for _, r := range v.Results {
					results = append(results, frame.newVar(r))
				}
				frame.push(Return{
					Results: results,
				})
				if len(state.frames) > 1 {
					state.frames = state.frames[:len(state.frames)-1]
					queue.push(state)
				} else {
					if model, sat := solve(fn, state.formula()); sat {
						fmt.Println("found solution for path:", state.frames[0].blockOrder)
						fmt.Println(model)
						testcases = append(testcases, Testcase{model: model})
					}
				}
				break instructionLoop
			case *ssa.UnOp:
				frame.push(UnOp{
					Result: frame.newVar(v),
					Arg:    frame.newVar(v.X),
					Op:     v.Op.String(),
				})
			case *ssa.Call:
				var args []Var
				for _, a := range v.Call.Args {
					args = append(args, frame.newVar(a))
				}
				name := removeArgs(v.Call.String())
				if IsBuiltIn(name) {
					frame.push(BuiltInCall{
						Result: frame.newVar(v),
						Name:   name,
						Args:   args,
					})
				} else {
					fn := v.Call.StaticCallee()
					if fn.Package().String() != pkg.String() {
						panic("external calls are not supported")
					}
					nextCall := &DynamicCall{
						Result: frame.newVar(v),
						Name:   name,
						Args:   args,
						Params: nil,
						Body:   nil,
					}
					frame.push(nextCall)
					frame.nextInstr = index + 1
					state.nextFrameId++
					nextFrame := &Frame{id: state.nextFrameId, function: fn, call: nextCall}
					state.frames = append(state.frames, nextFrame)
					for _, p := range fn.Params {
						tmp := &TempRegister{t: p.Type(), name: p.Name()}
						nextCall.Params = append(nextCall.Params, nextFrame.newVar(tmp))
					}
					if _, sat := solve(fn, state.formula()); sat {
						queue.push(state)
					}
					break instructionLoop
				}
			case *ssa.Convert:
				frame.push(Convert{
					Result: frame.newVar(v),
					Arg:    frame.newVar(v.X),
				})
			case *ssa.Phi:
				preds := v.Block().Preds
				var blocksIdxs []int
				for _, b := range preds {
					blocksIdxs = append(blocksIdxs, b.Index)
				}
				mostRecent := 0
				for _, i := range frame.blockOrder[:len(frame.blockOrder)-1] {
					for k, j := range blocksIdxs {
						if j == i {
							mostRecent = k
						}
					}
				}
				frame.push(Convert{
					Result: frame.newVar(v),
					Arg:    frame.newVar(v.Edges[mostRecent]),
				})
			case *ssa.IndexAddr:
				frame.push(IndexAddr{
					Result: frame.newVar(v),
					Array:  frame.newVar(v.X),
					Index:  frame.newVar(v.Index),
				})
			case *ssa.FieldAddr:
				frame.push(FieldAddr{
					Result: frame.newVar(v),
					Struct: frame.newVar(v.X),
					Field:  v.Field,
				})
			default:
				panic(fmt.Sprint("unknown instruction: '", v.String(), "'"))
			}
		}
	}
	return testcases
}

func solve(fn *ssa.Function, f Formula) (model *z3.Model, sat bool) {
	vars := make(map[string]Var, 0)
	f.ScanVars(vars)
	vars[resultSpecialVar] = Var{
		Name:     resultSpecialVar,
		Type:     fn.Signature.Results().At(0).Type(),
		Constant: false,
	}

	z3ctx := z3.NewContext(nil)
	ctx := &EncodingContext{
		Context: z3ctx,

		vars:     make(map[string]SymValue, 0),
		rawTypes: make(map[string]z3.Sort, 0),

		varsUsed: make(map[string]struct{}),
		varCount: make(map[string]int),

		fieldsMemory:      make(map[string][]z3.Array),
		valuesMemory:      make(map[string]z3.Array),
		arrayValuesMemory: make(map[string]z3.Array),
		arrayLenMemory:    make(map[string]z3.Array),

		floatSort:   z3ctx.FloatSort(11, 53),
		complexSort: z3ctx.UninterpretedSort("complex128"),
		stringSort:  z3ctx.UninterpretedSort("string"),

		addrSort: z3ctx.UninterpretedSort("$addr"),
	}

	for _, v := range vars {
		ctx.AddType(v.Type)
	}

	for _, v := range vars {
		ctx.AddVar(v.Name, v.Name, v.Type)
	}

	return solveWithTimeout(f.Encode(ctx).(z3.Bool), ctx)
}

func solveWithTimeout(f z3.Bool, ctx *EncodingContext) (model *z3.Model, sat bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[WARNING]", r)
		}
	}()

	ctx.Config().SetUint("timeout", 15*1000)

	solver := z3.NewSolver(ctx.Context)
	solver.Assert(f)
	for _, a := range ctx.asserts {
		solver.Assert(a)
	}

	sat, err := solver.Check()
	if err != nil {
		panic(err)
	}

	if sat {
		return solver.Model(), true
	} else {
		return nil, false
	}
}
