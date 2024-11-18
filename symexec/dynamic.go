package symexec

import (
	"fmt"
	"math/rand"
	"runtime/debug"

	"github.com/aclements/go-z3/z3"
	"golang.org/x/tools/go/ssa"
)

func AnalyzeFileDynamic(filename string) map[string]bool {
	main := buildPackage(filename)
	res := make(map[string]bool, 0)
	for _, v := range main.Members {
		if fn, ok := v.(*ssa.Function); ok && fn.Name() != "init" {
			res[fn.Name()] = dynamicFunction(fn)
		}
	}
	fmt.Println()
	return res
}

func dynamicFunction(fn *ssa.Function) (ok bool) {
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
	execute(fn)
	return true
}

type State struct {
	blockOrder []int
}

func execute(fn *ssa.Function) {
	queue := []State{
		{blockOrder: []int{0}},
	}
	for len(queue) > 0 {
		index := rand.Intn(len(queue))
		next := queue[index]
		queue = append(queue[:index], queue[index+1:]...)
		var subFormulas []Formula
		var lastInstr ssa.Instruction
		for blockNumber, blockIndex := range next.blockOrder {
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
					if blockNumber+1 < len(next.blockOrder) {
						isTrue := false
						if v.Block().Succs[0].Index == next.blockOrder[blockNumber+1] {
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
					for _, i := range next.blockOrder {
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
		formula := And{SubFormulas: subFormulas}
		if model, sat := solve(formula); sat {
			switch v := lastInstr.(type) {
			case *ssa.If:
				thenState := State{}
				thenState.blockOrder = append(thenState.blockOrder, next.blockOrder...)
				thenState.blockOrder = append(thenState.blockOrder, v.Block().Succs[0].Index)
				queue = append(queue, thenState)
				elseState := State{}
				elseState.blockOrder = append(elseState.blockOrder, next.blockOrder...)
				elseState.blockOrder = append(elseState.blockOrder, v.Block().Succs[1].Index)
				queue = append(queue, elseState)
			case *ssa.Jump:
				newState := State{}
				newState.blockOrder = append(newState.blockOrder, next.blockOrder...)
				newState.blockOrder = append(newState.blockOrder, v.Block().Succs[0].Index)
				queue = append(queue, newState)
			case *ssa.Return:
				fmt.Println("found solution for path:", next.blockOrder)
				fmt.Println(model)
			default:
				panic(fmt.Sprint("unknown divergence instruction: '", v.String(), "'"))
			}
		}
	}
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
		ctx.AddVar(v)
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
