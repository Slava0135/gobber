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
			vars := parseVars(tc.model)
			args, err := initArgs(fn, vars)
			if err != nil {
				fmt.Println("[ERROR]", err)
				continue
			}
			name := functionName(fn)
			var argsNames []string
			for _, param := range fn.Params {
				name := param.Name()
				argsNames = append(argsNames, name)
			}
			argsStr := strings.Join(argsNames, ", ")
			f.WriteString(fmt.Sprintf("func Test_%s_%d(t *testing.T) {\n", name, i+1))
			results := fn.Signature.Results()
			if results != nil && results.Len() == 1 {
				want, err := parseResult(results.At(0).Type(), vars)
				if err != nil {
					fmt.Println("[ERROR]", err)
					continue
				}
				for _, code := range args {
					f.WriteString(fmt.Sprintf("\t%s\n", strings.ReplaceAll(code, "\n", "\n\t")))
				}
				f.WriteString(fmt.Sprintf("\tgot := %s(%s)\n", name, argsStr))
				f.WriteString(fmt.Sprintf("\t%s\n", strings.ReplaceAll(want, "\n", "\n\t")))
				f.WriteString("\tif got != want {\n")
				f.WriteString(fmt.Sprintf("\t\tt.Errorf(\"%s(%s) = %%v; want %%v\", got, want)\n", name, argsStr))
				f.WriteString("\t}\n")
			} else {
				f.WriteString(fmt.Sprintf("\t%s(%s)\n", name, argsStr))
			}
			f.WriteString("}\n\n")
		}
	}
}

func functionName(fn *ssa.Function) string {
	segments := strings.Split(fn.Name(), ".")
	return segments[len(segments)-1]
}

func parseVars(model *z3.Model) map[string]string {
	vars := make(map[string]string)
	for _, line := range strings.Split(model.String(), "\n") {
		segments := strings.Split(line, " -> ")
		if len(segments) == 2 {
			vars[segments[0]] = segments[1]
		}
	}
	return vars
}

func initArgs(fn *ssa.Function, vars map[string]string) (map[string]string, error) {
	args := make(map[string]string)
	for _, param := range fn.Params {
		name := param.Name()
		value := vars[name]
		code, err := initValue(name, value, param.Type())
		if err != nil {
			return nil, err
		}
		args[name] = code
	}
	return args, nil
}

func parseResult(t types.Type, vars map[string]string) (string, error) {
	value, ok := vars[resultSpecialVar]
	if !ok {
		return "", fmt.Errorf("result not found in model")
	}
	return initValue("want", value, t)
}

func initValue(name string, value string, t types.Type) (string, error) {
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
	switch t := t.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Int:
			var goValue string
			if value == "" {
				goValue = "0"
			} else {
				i, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return "", fmt.Errorf("error when parsing integer '%s': %w", value, err)
				}
				goValue = fmt.Sprint(i)
			}
			return fmt.Sprintf("%s := %s", name, goValue), nil
		case types.Bool:
			var goValue string
			if value == "" {
				goValue = "false"
			} else {
				b, err := strconv.ParseBool(value)
				if err != nil {
					return "", fmt.Errorf("error when parsing boolean '%s': %w", value, err)
				}
				goValue = fmt.Sprint(b)
			}
			return fmt.Sprintf("%s := %s", name, goValue), nil
		default:
			return "", fmt.Errorf("unknown basic type '%s'", t)
		}
	default:
		return "", fmt.Errorf("unknown type '%s'", t)
	}
}

func parseSmtFloat(str string) string {
	return ""
}
