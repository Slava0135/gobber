package symexec

import (
	"fmt"

	"golang.org/x/tools/go/ssa"
)

func AnalyzeFileDynamic(filename string) map[string]bool {
	main := buildPackage(filename)
	res := make(map[string]bool, 0)
	for _, v := range main.Members {
		if fn, ok := v.(*ssa.Function); ok && fn.Name() != "init" {
			res[fn.Name()] = false
		}
	}
	fmt.Println()
	return res
}
