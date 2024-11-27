package symexec

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type Register interface {
	Type() types.Type
	Name() string
}

func buildPackage(filename string) *ssa.Package {
	fmt.Printf(":: building SSA graph for file '%s'\n", filename)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		panic(err)
	}

	files := []*ast.File{f}

	pkg := types.NewPackage("main", "")

	main, _, err := ssautil.BuildPackage(&types.Config{Importer: importer.Default()}, fset, pkg, files, 0)
	if err != nil {
		panic(err)
	}
	return main
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
