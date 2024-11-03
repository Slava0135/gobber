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
				fmt.Println("::", "analyzing function", "'"+fn.Name()+"'")
				fmt.Println("::", "printing SSA blocks")
				printBlocks(fn)
				fmt.Println("::", "building formula")
				makeFormula(fn)
			}
		}

		fmt.Println()
	}
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

func makeFormula(fn *ssa.Function) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	f := getBlockFormula(fn.Blocks, 0, make([]bool, len(fn.Blocks)))
	fmt.Println(f)
}

func getBlockFormula(blocks []*ssa.BasicBlock, blockIndex int, visited []bool) Formula {
	if visited[blockIndex] {
		panic("[ERROR] cycles are not supported!")
	}
	visited[blockIndex] = true
	block := blocks[blockIndex]
	var subFormulas []Formula
	for _, v := range block.Instrs {
		switch v := v.(type) {
		case *ssa.BinOp:
			subFormulas = append(subFormulas, BinOp{
				Result: Var{v.Name()},
				Left:   Var{v.X.Name()},
				Op:     Op{v.Op.String()},
				Right:  Var{v.Y.Name()},
			})
		case *ssa.If:
			subFormulas = append(subFormulas, If{
				Cond: Var{Name: v.Cond.Name()},
				Then: getBlockFormula(blocks, block.Succs[0].Index, visited),
				Else: getBlockFormula(blocks, block.Succs[1].Index, visited),
			})
		case *ssa.Jump:
			subFormulas = append(subFormulas, getBlockFormula(blocks, block.Succs[0].Index, visited))
		case *ssa.Return:
			var results []Var
			for _, r := range v.Results {
				results = append(results, Var{r.Name()})
			}
			subFormulas = append(subFormulas, Return{
				Results: results,
			})
		case *ssa.UnOp:
			subFormulas = append(subFormulas, UnOp{
				Result: Var{v.Name()},
				Arg:    Var{v.X.Name()},
				Op:     Op{v.Op.String()},
			})
		default:
			panic(fmt.Sprint("[ERROR] unknown instruction: '", v.String(), "'"))
		}
	}
	return And{SubFormulas: subFormulas}
}
