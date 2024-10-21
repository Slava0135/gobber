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

func SSA() {
	os.Chdir("testdata")

	testcases, err := os.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, tc := range testcases {
		fmt.Printf(":: building ssa graph for file '%s'\n", tc.Name())

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
				fmt.Println("::", fn.Name())
				for _, v := range fn.Blocks {
					fmt.Println(v.String(), "->")
					for _, v := range v.Instrs {
						printInstr := func (name string) {
							fmt.Printf("  [%10s] %s\n", strings.ToUpper(name), v.String())
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
		}

		fmt.Println()
	}
}
