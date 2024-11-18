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

type Register interface {
	Type() types.Type
	Name() string
}

func Static() {
	os.Chdir("testdata")

	testcases, err := os.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, tc := range testcases {
		AnalyzeFile(tc.Name())
	}
}

func AnalyzeFile(filename string) map[string]bool {
	main := buildPackage(filename)
	res := make(map[string]bool, 0)
	for _, v := range main.Members {
		if fn, ok := v.(*ssa.Function); ok && fn.Name() != "init" {
			res[fn.Name()] = staticFunction(fn)
		}
	}
	fmt.Println()
	return res
}

func staticFunction(fn *ssa.Function) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[ERROR]", r)
			fmt.Println(string(debug.Stack()))
		}
	}()
	fmt.Println("::", "analyzing function", "'"+fn.Name()+"'")
	fmt.Println("::", "printing SSA blocks")
	printBlocks(fn)
	fmt.Println("::", "building formula")
	f := makeFormula(fn)
	fmt.Println("::", "encoding formula")
	encodeFormula(fn, f)
	return true
}

func makeFormula(fn *ssa.Function) Formula {
	f := getBlockFormula(fn.Blocks, 0, make([]int, len(fn.Blocks)), 1)
	fmt.Println("::", "logical")
	fmt.Println(f)
	// fmt.Println("::", "yaml")
	// fmt.Println(toYaml(f))
	return f
}

func getBlockFormula(blocks []*ssa.BasicBlock, blockIndex int, visitOrder []int, depth int) Formula {
	if visitOrder[blockIndex] > 0 {
		panic("cycles are not supported!")
	}

	newVisitOrder := make([]int, len(visitOrder))
	copy(newVisitOrder, visitOrder)
	newVisitOrder[blockIndex] = depth

	block := blocks[blockIndex]
	var subFormulas []Formula
	for _, v := range block.Instrs {
		switch v := v.(type) {
		case *ssa.BinOp:
			subFormulas = append(subFormulas, BinOp{
				Result: NewVar(v),
				Left:   NewVar(v.X),
				Op:     v.Op.String(),
				Right:  NewVar(v.Y),
			})
		case *ssa.If:
			subFormulas = append(subFormulas, If{
				Cond: NewVar(v.Cond),
				Then: getBlockFormula(blocks, block.Succs[0].Index, newVisitOrder, depth+1),
				Else: getBlockFormula(blocks, block.Succs[1].Index, newVisitOrder, depth+1),
			})
		case *ssa.Jump:
			subFormulas = append(subFormulas, getBlockFormula(blocks, block.Succs[0].Index, newVisitOrder, depth+1))
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
			mostRecent := 0
			preds := v.Block().Preds
			for i, b := range preds {
				if visitOrder[b.Index] > visitOrder[preds[mostRecent].Index] {
					mostRecent = i
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
	return And{SubFormulas: subFormulas}
}

func encodeFormula(fn *ssa.Function, f Formula) {
	fmt.Println("::", "listing all variables")
	vars := make(map[string]Var, 0)
	f.ScanVars(vars)
	for _, v := range vars {
		fmt.Print(v, " ")
	}
	fmt.Println()

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

	fmt.Println("::", "encoding formula in Z3")
	encodedFormula := f.Encode(ctx).(z3.Bool)
	fmt.Println(encodedFormula)

	fmt.Println("::", "solving")
	solver := z3.NewSolver(ctx.Context)
	solver.Assert(encodedFormula)
	for _, a := range ctx.asserts {
		solver.Assert(a)
	}
	sat, err := solver.Check()
	if err != nil {
		panic(err)
	}
	if !sat {
		panic("unexpected unsat")
	}
	fmt.Println("SAT")
	fmt.Println(strings.TrimSpace(solver.Model().String()))
}
