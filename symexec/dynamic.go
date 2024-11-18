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
		for _, i := range next.blockOrder {
			block := fn.Blocks[i]
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
					// TODO
				case *ssa.Jump:
					// TODO
				case *ssa.Return:
					// TODO
				case *ssa.UnOp:
					subFormulas = append(subFormulas, UnOp{
						Result: NewVar(v),
						Arg:    NewVar(v.X),
						Op:     v.Op.String(),
					})
				case *ssa.Call:
					// TODO
				case *ssa.Convert:
					subFormulas = append(subFormulas, Convert{
						Result: NewVar(v),
						Arg:    NewVar(v.X),
					})
				case *ssa.Phi:
					// TODO
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
		if model, ok := solve(formula); ok {
			switch v := lastInstr.(type) {
			case *ssa.If:
			case *ssa.Jump:
			case *ssa.Return:
				fmt.Println("found solution for path:", next.blockOrder)
				fmt.Println(model)
			default:
				panic(fmt.Sprint("unknown divergence instruction: '", v.String(), "'"))
			}
		}
	}
}

func solve(formula Formula) (model z3.Model, ok bool) {
	return z3.Model{}, true
}
