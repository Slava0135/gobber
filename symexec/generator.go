package symexec

import (
	"fmt"
	"go/types"
	"math"
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

func trim(value string) string {
	var trimmed []rune
	for _, c := range value {
		switch c {
		case '(', ')', '\n', '\t', ' ':
			continue
		default:
			trimmed = append(trimmed, c)
		}
	}
	return string(trimmed)
}

func initValue(name string, value string, t types.Type) (string, error) {
	switch t := t.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Int:
			value := trim(value)
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
			value := trim(value)
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
		case types.Float64:
			if value == "" {
				return fmt.Sprintf("%s := 0.0", name), nil
			} else {
				return initSmtFloat64(name, value)
			}
		default:
			return "", fmt.Errorf("unknown basic type '%s'", t)
		}
	default:
		return "", fmt.Errorf("unknown type '%s'", t)
	}
}

func initSmtFloat64(name string, value string) (string, error) {
	value = strings.Trim(value, "()")
	components := strings.Split(value, " ")
	if len(components) != 4 {
		return "", fmt.Errorf("expected 4 components for float64: %s", value)
	}
	if components[0] != "_" {
		signBin, ok := strings.CutPrefix(components[1], "#b")
		if !ok {
			return "", fmt.Errorf("invalid sign for float64: %s", value)
		}
		expBin, ok := strings.CutPrefix(components[2], "#b")
		if !ok {
			return "", fmt.Errorf("invalid exponent for float64: %s", value)
		}
		mantHex, ok := strings.CutPrefix(components[3], "#x")
		if !ok {
			return "", fmt.Errorf("invalid mantissa for float64: %s", value)
		}
		signExpBin := signBin + expBin
		hex := ""
		for i := 0; i < len(signExpBin); i += 4 {
			v, err := strconv.ParseUint(signExpBin[i:i+4], 2, 4)
			if err != nil {
				return "", fmt.Errorf("error when parsing float64 '%s': %w", value, err)
			}
			hex += fmt.Sprintf("%x", v)
		}
		hex += mantHex
		bits, err := strconv.ParseUint(hex, 16, 64)
		if err != nil {
			return "", fmt.Errorf("error when parsing float64 '%s': %w", value, err)
		}
		f64 := math.Float64frombits(bits)
		return fmt.Sprintf("%s := %f", name, f64), nil
	} else {
		return "???", nil
	}
}
