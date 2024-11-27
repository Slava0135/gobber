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
		AnalyzeFileDynamic(tc.Name())
	}
}

func AnalyzeFileDynamic(filename string) map[string]bool {
	main := buildPackage(filename)
	res := make(map[string]bool, 0)
	for _, v := range main.Members {
		if fn, ok := v.(*ssa.Function); ok && fn.Name() != "init" {
			res[fn.Name()] = dynamicFunction(fn, main)
		}
	}
	fmt.Println()
	return res
}

func dynamicFunction(fn *ssa.Function, pkg *ssa.Package) (ok bool) {
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
	execute(fn, pkg, &RandomQueue{})
	return true
}

type State struct {
	frames []*Frame
}

type Frame struct {
	function   *ssa.Function
	blockOrder []int
}

func execute(fn *ssa.Function, pkg *ssa.Package, queue Queue) {
	startFrame := &Frame{function: fn, blockOrder: []int{0}}
	queue.push(State{frames: []*Frame{startFrame}})
	for !queue.empty() {
		next := queue.pop()
		var subFormulas []Formula
		var lastInstr ssa.Instruction
		for _, frame := range next.frames {
			for blockNumber, blockIndex := range frame.blockOrder {
				block := fn.Blocks[blockIndex]
				for _, v := range block.Instrs {
					lastInstr = v
					switch v := v.(type) {
					case *ssa.BinOp:
						subFormulas = append(subFormulas, BinOp{
							Result: NewVar(v),
							Left:   NewVar(v.X),
							Op:     v.Op.String(),
							Right:  NewVar(v.Y),
						})
					case *ssa.If:
						if blockNumber+1 < len(frame.blockOrder) {
							isTrue := false
							if v.Block().Succs[0].Index == frame.blockOrder[blockNumber+1] {
								isTrue = true
							}
							subFormulas = append(subFormulas, Condition{
								Cond:   NewVar(v.Cond),
								IsTrue: isTrue,
							})
						}
					case *ssa.Jump:
						// do nothing
					case *ssa.Return:
						var results []Var
						for _, r := range v.Results {
							results = append(results, NewVar(r))
						}
						subFormulas = append(subFormulas, Return{
							Results: results,
						})
					case *ssa.UnOp:
						subFormulas = append(subFormulas, UnOp{
							Result: NewVar(v),
							Arg:    NewVar(v.X),
							Op:     v.Op.String(),
						})
					case *ssa.Call:
						// TODO make interprocedural
						var args []Var
						for _, a := range v.Call.Args {
							args = append(args, NewVar(a))
						}
						subFormulas = append(subFormulas, Call{
							Result: NewVar(v),
							Name:   removeArgs(v.Call.String()),
							Args:   args,
						})
					case *ssa.Convert:
						subFormulas = append(subFormulas, Convert{
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
						for _, i := range frame.blockOrder[:blockNumber] {
							for k, j := range blocksIdxs {
								if j == i {
									mostRecent = k
								}
							}
						}
						subFormulas = append(subFormulas, Convert{
							Result: NewVar(v),
							Arg:    NewVar(v.Edges[mostRecent]),
						})
					case *ssa.IndexAddr:
						subFormulas = append(subFormulas, IndexAddr{
							Result: NewVar(v),
							Array:  NewVar(v.X),
							Index:  NewVar(v.Index),
						})
					case *ssa.FieldAddr:
						subFormulas = append(subFormulas, FieldAddr{
							Result: NewVar(v),
							Struct: NewVar(v.X),
							Field:  v.Field,
						})
					default:
						panic(fmt.Sprint("unknown instruction: '", v.String(), "'"))
					}
				}
			}
		}
		formula := And{SubFormulas: subFormulas}
		if model, sat := solve(formula); sat {
			switch v := lastInstr.(type) {
			case *ssa.If:
				{
					thenState := next.copy()
					lastFrame := thenState.frames[len(thenState.frames)-1]
					lastFrame.blockOrder = append(lastFrame.blockOrder, v.Block().Succs[0].Index)
					queue.push(thenState)
				}
				{
					elseState := next.copy()
					lastFrame := elseState.frames[len(elseState.frames)-1]
					lastFrame.blockOrder = append(lastFrame.blockOrder, v.Block().Succs[1].Index)
					queue.push(elseState)
				}
			case *ssa.Jump:
				jumpState := next.copy()
				lastFrame := jumpState.frames[len(jumpState.frames)-1]
				lastFrame.blockOrder = append(lastFrame.blockOrder, v.Block().Succs[0].Index)
				queue.push(jumpState)
			case *ssa.Return:
				fmt.Println("found solution for path:", next.frames[0].blockOrder)
				fmt.Println(model)
			default:
				panic(fmt.Sprint("unknown divergence instruction: '", v.String(), "'"))
			}
		}
	}
}

func (s *State) copy() State {
	stateCopy := State{}
	for _, frame := range s.frames {
		frameCopy := &Frame{}
		frameCopy.function = frame.function
		frameCopy.blockOrder = append(frameCopy.blockOrder, frame.blockOrder...)
		stateCopy.frames = append(stateCopy.frames, frameCopy)
	}
	return stateCopy
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
