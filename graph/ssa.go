package graph

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"strings"

	"github.com/aclements/go-z3/z3"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type Register interface {
	Type() types.Type
	Name() string
}

func SSA() {
	os.Chdir("testdata")

	testcases, err := os.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, tc := range testcases {
		fmt.Printf(":: building SSA graph for file '%s'\n", tc.Name())

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, tc.Name(), nil, 0)
		if err != nil {
			panic(err)
		}

		files := []*ast.File{f}

		pkg := types.NewPackage("main", "")

		main, _, err := ssautil.BuildPackage(&types.Config{Importer: importer.Default()}, fset, pkg, files, 0)
		if err != nil {
			panic(err)
		}

		for _, v := range main.Members {
			if fn, ok := v.(*ssa.Function); ok && fn.Name() != "init" {
				doSSA(fn)
			}
		}

		fmt.Println()
	}
}

func doSSA(fn *ssa.Function) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[ERROR]", r)
		}
	}()
	fmt.Println("::", "analyzing function", "'"+fn.Name()+"'")
	fmt.Println("::", "printing SSA blocks")
	printBlocks(fn)
	fmt.Println("::", "building formula")
	f := makeFormula(fn)
	fmt.Println("::", "encoding formula")
	encodeFormula(fn, f)
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
				Result: Var{Name: removeType(v.Name()), Type: v.Type().String()},
				Left:   Var{Name: removeType(v.X.Name()), Type: v.X.Type().String(), Constant: isConstant(v.X.Name())},
				Op:     v.Op.String(),
				Right:  Var{Name: removeType(v.Y.Name()), Type: v.Y.Type().String(), Constant: isConstant(v.Y.Name())},
			})
		case *ssa.If:
			subFormulas = append(subFormulas, If{
				Cond: Var{Name: v.Cond.Name(), Type: v.Cond.Type().String(), Constant: isConstant(v.Cond.Name())},
				Then: getBlockFormula(blocks, block.Succs[0].Index, newVisitOrder, depth+1),
				Else: getBlockFormula(blocks, block.Succs[1].Index, newVisitOrder, depth+1),
			})
		case *ssa.Jump:
			subFormulas = append(subFormulas, getBlockFormula(blocks, block.Succs[0].Index, newVisitOrder, depth+1))
		case *ssa.Return:
			var results []Var
			for _, r := range v.Results {
				results = append(results, Var{Name: removeType(r.Name()), Type: r.Type().String(), Constant: isConstant(r.Name())})
			}
			subFormulas = append(subFormulas, Return{
				Results: results,
			})
		case *ssa.UnOp:
			subFormulas = append(subFormulas, UnOp{
				Result: Var{Name: removeType(v.Name()), Type: v.Type().String()},
				Arg:    Var{Name: removeType(v.X.Name()), Type: v.X.Type().String(), Constant: isConstant(v.X.Name())},
				Op:     v.Op.String(),
			})
		case *ssa.Call:
			var args []Var
			for _, a := range v.Call.Args {
				args = append(args, Var{Name: removeType(a.Name()), Type: a.Type().String(), Constant: isConstant(a.Name())})
			}
			subFormulas = append(subFormulas, Call{
				Result: Var{Name: v.Name(), Type: v.Type().String()},
				Name:   removeArgs(v.Call.String()),
				Args:   args,
			})
		case *ssa.Convert:
			subFormulas = append(subFormulas, Convert{
				Result: Var{Name: v.Name(), Type: v.Type().String()},
				Arg:    Var{Name: v.X.Name(), Type: v.X.Type().String(), Constant: isConstant(v.X.Name())},
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
				Result: Var{Name: v.Name(), Type: v.Type().String()},
				Arg:    Var{Name: v.Edges[mostRecent].Name(), Type: v.Edges[mostRecent].Type().String()},
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

	ctx := z3.NewContext(nil)

	encodedVars := make(map[string]z3.Value, 0)
	for _, v := range vars {
		switch v.Type {
		case intType:
			encodedVars[v.Name] = ctx.IntConst(v.Name)
		case boolType:
			encodedVars[v.Name] = ctx.BoolConst(v.Name)
		default:
			panic(fmt.Sprintf("unknown type '%s'", v.Type))
		}
	}
	funcs := make(map[string]z3.FuncDecl, 0)

	encodedFormula := f.Encode(encodedVars, funcs)
	fmt.Println(encodedFormula)
}
