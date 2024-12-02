package symexec

import (
	"fmt"
	"go/types"
	"os"
	"strconv"
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
			args, err := parseArgs(fn, tc.model)
			if err != nil {
				fmt.Println("[ERROR]", err)
				continue
			}
			argsStr := strings.Join(args, ", ")
			name := functionName(fn)
			f.WriteString(fmt.Sprintf("func Test_%s_%d(t *testing.T) {\n", name, i+1))
			f.WriteString(fmt.Sprintf("\t%s(%s)\n", name, argsStr))
			f.WriteString("}\n\n")
		}
	}
}

func functionName(fn *ssa.Function) string {
	segments := strings.Split(fn.Name(), ".")
	return segments[len(segments)-1]
}

func parseArgs(fn *ssa.Function, model *z3.Model) ([]string, error) {
	var args []string
	vars := make(map[string]string)
	for _, line := range strings.Split(model.String(), "\n") {
		segments := strings.Split(line, " -> ")
		if len(segments) == 2 {
			vars[segments[0]] = segments[1]
		}
	}
	for _, param := range fn.Params {
		name := param.Name()
		value, ok := vars[name]
		if !ok {
			return nil, fmt.Errorf("param named '%s' not found in model", name)
		}
		var trimmed []rune
		for _, c := range value {
			switch c {
			case '(', ')', '\n', '\t', ' ':
				continue
			default:
				trimmed = append(trimmed, c)
			}
		}
		value = string(trimmed)
		var goValue string
		switch t := param.Type().(type) {
		case *types.Basic:
			switch t.Kind() {
			case types.Int:
				i, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("error when parsing integer '%s': %w", value, err)
				}
				goValue = fmt.Sprint(i)
			case types.Bool:
				b, err := strconv.ParseBool(value)
				if err != nil {
					return nil, fmt.Errorf("error when parsing boolean '%s': %w", value, err)
				}
				goValue = fmt.Sprint(b)
			default:
				return nil, fmt.Errorf("unknown basic type '%s'", param.Type())
			}
		default:
			return nil, fmt.Errorf("unknown type '%s'", param.Type())
		}
		args = append(args, goValue)
	}
	return args, nil
}
