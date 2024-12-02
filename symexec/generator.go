package symexec

import (
	"fmt"
	"os"
	"strings"

	"github.com/aclements/go-z3/z3"
	"golang.org/x/tools/go/ssa"
)

type Testcase struct {
	model *z3.Model
}

func GenerateTests(filename string, functionTestcases map[*ssa.Function][]Testcase) {
	filenameWithoutExt, _ := strings.CutSuffix(filename, ".go")
	fmt.Println(":: generating tests")
	f, err := os.Create(filenameWithoutExt + "_test.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	prelude := `
package main

import "testing"
`
	f.WriteString(strings.Trim(prelude, "\n"))
	f.WriteString("\n\n")
	for fn, testcases := range functionTestcases {
		for i, tc := range testcases {
			name := functionName(fn)
			f.WriteString(fmt.Sprintf("func Test_%s_%d(t *testing.T) {\n", name, i+1))
			args := strings.Join(parseArgs(fn, tc.model), ", ")
			f.WriteString(fmt.Sprintf("\t%s(%s)\n", name, args))
			f.WriteString("}\n\n")
		}
	}
}

func functionName(fn *ssa.Function) string {
	segments := strings.Split(fn.Name(), ".")
	return segments[len(segments)-1]
}

func parseArgs(fn *ssa.Function, model *z3.Model) []string {
	var args []string
	return args
}
