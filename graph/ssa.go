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
						printInstrReg := func (reg string, name string) {
							fmt.Printf("  [%10s] %s <- %s\n", strings.ToUpper(name), reg, v.String())
						}
						switch v := v.(type) {
						case *ssa.Alloc:
							printInstrReg(v.Name(), "alloc")
						case *ssa.BinOp:
							printInstrReg(v.Name(), "binop")
						case *ssa.Call:
							printInstrReg(v.Name(), "call")
						case *ssa.Convert:
							printInstrReg(v.Name(), "convert")
						case *ssa.Extract:
							printInstrReg(v.Name(), "extract")
						case *ssa.Field:
							printInstrReg(v.Name(), "field")
						case *ssa.FieldAddr:
							printInstrReg(v.Name(), "field addr")
						case *ssa.If:
							printInstr("if")
						case *ssa.Index:
							printInstrReg(v.Name(), "index")
						case *ssa.IndexAddr:
							printInstrReg(v.Name(), "index addr")
						case *ssa.Jump:
							printInstr("jump")
						case *ssa.Lookup:
							printInstrReg(v.Name(), "lookup")
						case *ssa.MakeMap:
							printInstrReg(v.Name(), "make map")
						case *ssa.MakeSlice:
							printInstrReg(v.Name(), "make slice")
						case *ssa.MapUpdate:
							printInstr("map update")
						case *ssa.Phi:
							printInstrReg(v.Name(), "phi")
						case *ssa.Return:
							printInstr("return")
						case *ssa.Select:
							printInstrReg(v.Name(), "select")
						case *ssa.Store:
							printInstr("store")
						case *ssa.UnOp:
							printInstrReg(v.Name(), "unop")
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
