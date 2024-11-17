package symexec

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"runtime/debug"
	"strings"

	"github.com/aclements/go-z3/z3"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
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

func AnalyzeFile(tc string) map[string]bool {
	fmt.Printf(":: building SSA graph for file '%s'\n", tc)
	
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, tc, nil, 0)
	if err != nil {
		panic(err)
	}
	
	files := []*ast.File{f}
	
	pkg := types.NewPackage("main", "")
	
	main, _, err := ssautil.BuildPackage(&types.Config{Importer: importer.Default()}, fset, pkg, files, 0)
	if err != nil {
		panic(err)
	}
	
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

func printBlocks(fn *ssa.Function) {
	for _, v := range fn.Blocks {
		fmt.Println(v.String(), "->")
		for _, v := range v.Instrs {
			printInstr := func(name string) {
				if reg, ok := v.(Register); ok {
					fmt.Printf("  [%10s] %s:%s <-- %s\n", strings.ToUpper(name), reg.Name(), reg.Type(), v.String())
				} else {
					fmt.Printf("  [%10s] %s\n", strings.ToUpper(name), v.String())
				}
			}
			switch v.(type) {
			case *ssa.Alloc:
				printInstr("alloc")
			case *ssa.BinOp:
				printInstr("binop")
			case *ssa.Call:
				printInstr("call")
			case *ssa.Convert:
				printInstr("convert")
			case *ssa.Extract:
				printInstr("extract")
			case *ssa.Field:
				printInstr("field")
			case *ssa.FieldAddr:
				printInstr("field addr")
			case *ssa.If:
				printInstr("if")
			case *ssa.Index:
				printInstr("index")
			case *ssa.IndexAddr:
				printInstr("index addr")
			case *ssa.Jump:
				printInstr("jump")
			case *ssa.Lookup:
				printInstr("lookup")
			case *ssa.MakeMap:
				printInstr("make map")
			case *ssa.MakeSlice:
				printInstr("make slice")
			case *ssa.MapUpdate:
				printInstr("map update")
			case *ssa.Phi:
				printInstr("phi")
			case *ssa.Return:
				printInstr("return")
			case *ssa.Select:
				printInstr("select")
			case *ssa.Store:
				printInstr("store")
			case *ssa.UnOp:
				printInstr("unop")
			default:
				panic("unknown instruction")
			}
		}
	}
}

func makeFormula(fn *ssa.Function) Formula {
	f := getBlockFormula(fn.Blocks, 0, make([]int, len(fn.Blocks)), 1)
	fmt.Println("::", "logical")
	fmt.Println(f)
	// fmt.Println("::", "yaml")
	// fmt.Println(toYaml(f))
	return f
}

func removeType(str string) string {
	return strings.Split(str, ":")[0]
}

func isConstant(str string) bool {
	return strings.Contains(str, ":")
}

func removeArgs(str string) string {
	return strings.Split(str, "(")[0]
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
				Result: Var{Name: removeType(v.Name()), Type: v.Type()},
				Left:   Var{Name: removeType(v.X.Name()), Type: v.X.Type(), Constant: isConstant(v.X.Name())},
				Op:     v.Op.String(),
				Right:  Var{Name: removeType(v.Y.Name()), Type: v.Y.Type(), Constant: isConstant(v.Y.Name())},
			})
		case *ssa.If:
			subFormulas = append(subFormulas, If{
				Cond: Var{Name: v.Cond.Name(), Type: v.Cond.Type(), Constant: isConstant(v.Cond.Name())},
				Then: getBlockFormula(blocks, block.Succs[0].Index, newVisitOrder, depth+1),
				Else: getBlockFormula(blocks, block.Succs[1].Index, newVisitOrder, depth+1),
			})
		case *ssa.Jump:
			subFormulas = append(subFormulas, getBlockFormula(blocks, block.Succs[0].Index, newVisitOrder, depth+1))
		case *ssa.Return:
			var results []Var
			for _, r := range v.Results {
				results = append(results, Var{Name: removeType(r.Name()), Type: r.Type(), Constant: isConstant(r.Name())})
			}
			subFormulas = append(subFormulas, Return{
				Results: results,
			})
		case *ssa.UnOp:
			subFormulas = append(subFormulas, UnOp{
				Result: Var{Name: removeType(v.Name()), Type: v.Type()},
				Arg:    Var{Name: removeType(v.X.Name()), Type: v.X.Type(), Constant: isConstant(v.X.Name())},
				Op:     v.Op.String(),
			})
		case *ssa.Call:
			var args []Var
			for _, a := range v.Call.Args {
				args = append(args, Var{Name: removeType(a.Name()), Type: a.Type(), Constant: isConstant(a.Name())})
			}
			subFormulas = append(subFormulas, Call{
				Result: Var{Name: v.Name(), Type: v.Type()},
				Name:   removeArgs(v.Call.String()),
				Args:   args,
			})
		case *ssa.Convert:
			subFormulas = append(subFormulas, Convert{
				Result: Var{Name: v.Name(), Type: v.Type()},
				Arg:    Var{Name: removeType(v.X.Name()), Type: v.X.Type(), Constant: isConstant(v.X.Name())},
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
				Result: Var{Name: v.Name(), Type: v.Type()},
				Arg:    Var{Name: v.Edges[mostRecent].Name(), Type: v.Edges[mostRecent].Type()},
			})
		case *ssa.IndexAddr:
			subFormulas = append(subFormulas, IndexAddr{
				Result: Var{Name: v.Name(), Type: v.Type()},
				Array:  Var{Name: v.X.Name(), Type: v.X.Type()},
				Index:  Var{Name: removeType(v.Index.Name()), Type: v.Index.Type(), Constant: isConstant(v.Index.Name())},
			})
		case *ssa.FieldAddr:
			subFormulas = append(subFormulas, FieldAddr{
				Result: Var{Name: v.Name(), Type: v.Type()},
				Struct: Var{Name: v.X.Name(), Type: v.X.Type()},
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
