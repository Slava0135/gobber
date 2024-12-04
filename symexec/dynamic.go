package symexec

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/aclements/go-z3/z3"
	"golang.org/x/tools/go/ssa"
)

func Dynamic() {
	os.Chdir("testdata")

	testcases, err := os.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, tc := range testcases {
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
	frames []*Frame
}

func (s *State) copy() *State {
	stateCopy := &State{}
	for _, frame := range s.frames {
		stateCopy.frames = append(stateCopy.frames, frame.copy())
	}
	return stateCopy
}

func (s *State) currentFrame() *Frame {
	return s.frames[len(s.frames)-1]
}

func (s *State) formula() Formula {
	return And{s.frames[0].call.Body}
}

type Frame struct {
	function   *ssa.Function
	blockOrder []int
	call       *DynamicCall
	nextBlock  int
	nextInstr  int
}

func (frame *Frame) push(f Formula) {
	frame.call.Body = append(frame.call.Body, f)
}

func (frame *Frame) copy() *Frame {
	var blockOrder []int
	var body []Formula
	return &Frame{
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
		frame := state.currentFrame()
		block := frame.function.Blocks[frame.nextBlock]
		frame.blockOrder = append(frame.blockOrder, frame.nextBlock)
		for index, instr := range block.Instrs {
			if index < frame.nextInstr {
				continue
			}
			switch v := instr.(type) {
			case *ssa.BinOp:
				frame.push(BinOp{
					Result: NewVar(v),
					Left:   NewVar(v.X),
					Op:     v.Op.String(),
					Right:  NewVar(v.Y),
				})
			case *ssa.If:
				{
					thenState := state.copy()
					thenFrame := thenState.currentFrame()
					thenFrame.nextBlock = v.Block().Succs[0].Index
					thenFrame.push(Condition{
						Cond:   NewVar(v.Cond),
						IsTrue: true,
					})
					if _, sat := solve(thenState.formula()); sat {
						queue.push(thenState)
					}
				}
				{
					elseState := state.copy()
					elseFrame := elseState.currentFrame()
					elseFrame.nextBlock = v.Block().Succs[1].Index
					elseFrame.push(Condition{
						Cond:   NewVar(v.Cond),
						IsTrue: false,
					})
					if _, sat := solve(elseState.formula()); sat {
						queue.push(elseState)
					}
				}
			case *ssa.Jump:
				state.currentFrame().nextBlock = v.Block().Succs[0].Index
				queue.push(state)
			case *ssa.Return:
				var results []Var
				for _, r := range v.Results {
					results = append(results, NewVar(r))
				}
				frame.push(Return{
					Results: results,
				})
				if len(state.frames) > 1 {
					state.frames = state.frames[:len(state.frames)-1]
					queue.push(state)
				} else {
					if model, sat := solve(state.formula()); sat {
						fmt.Println("found solution for path:", state.frames[0].blockOrder)
						fmt.Println(model)
						testcases = append(testcases, Testcase{model: model})
					}
				}
			case *ssa.UnOp:
				frame.push(UnOp{
					Result: NewVar(v),
					Arg:    NewVar(v.X),
					Op:     v.Op.String(),
				})
			case *ssa.Call:
				var args []Var
				for _, a := range v.Call.Args {
					args = append(args, NewVar(a))
				}
				name := removeArgs(v.Call.String())
				if IsBuiltIn(name) {
					frame.push(BuiltInCall{
						Result: NewVar(v),
						Name:   name,
						Args:   args,
					})
				} else {
					var params []Var
					for _, p := range v.Common().StaticCallee().Params {
						tmp := &TempRegister{t: p.Type(), name: p.Name()}
						params = append(params, NewVar(tmp))
					}
					nextCall := &DynamicCall{
						Result: NewVar(v),
						Name:   name,
						Args:   args,
						Params: params,
						Body:   nil,
					}
					frame.push(nextCall)
					frame.nextInstr = index + 1
					nextFrame := &Frame{function: v.Call.StaticCallee(), blockOrder: []int{0}, call: nextCall}
					state.frames = append(state.frames, nextFrame)
					if _, sat := solve(state.formula()); sat {
						queue.push(state)
					}
				}
			case *ssa.Convert:
				frame.push(Convert{
					Result: NewVar(v),
					Arg:    NewVar(v.X),
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
					Result: NewVar(v),
					Arg:    NewVar(v.Edges[mostRecent]),
				})
			case *ssa.IndexAddr:
				frame.push(IndexAddr{
					Result: NewVar(v),
					Array:  NewVar(v.X),
					Index:  NewVar(v.Index),
				})
			case *ssa.FieldAddr:
				frame.push(FieldAddr{
					Result: NewVar(v),
					Struct: NewVar(v.X),
					Field:  v.Field,
				})
			default:
				panic(fmt.Sprint("unknown instruction: '", v.String(), "'"))
			}
		}
	}
	return testcases
}

func solve(f Formula) (model *z3.Model, sat bool) {
	vars := make(map[string]Var, 0)
	f.ScanVars(vars)

	z3ctx := z3.NewContext(nil)
	ctx := &EncodingContext{
		Context: z3ctx,

		vars:     make(map[string]SymValue, 0),
		funcs:    make(map[string]z3.FuncDecl, 0),
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

	encodedFormula := f.Encode(ctx).(z3.Bool)

	solver := z3.NewSolver(ctx.Context)
	solver.Assert(encodedFormula)
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
