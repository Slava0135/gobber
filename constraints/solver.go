package constraints

import (
	"fmt"
	"strings"

	"github.com/aclements/go-z3/z3"
)

const intSize = 64

func solve(solver *z3.Solver, path string, asserts ...z3.Bool) {
	printPath(path)
	for _, v := range asserts {
		solver.Assert(v)
	}
	sat, err := solver.Check()
	if err != nil {
		panic(err)
	}
	if !sat {
		panic("unexpected unsat")
	}
	fmt.Println(solver.Model())
	solver.Reset()
}

func solveIncrement(solver *z3.Solver, path string, asserts ...z3.Bool) {
	printPath(path)
	solver.Push()
	for _, v := range asserts {
		solver.Assert(v)
	}
	sat, err := solver.Check()
	if err != nil {
		panic(err)
	}
	if !sat {
		panic("unexpected unsat")
	}
	fmt.Println(solver.Model())
	solver.Pop()
}

func printSrc(src string) {
	maxLen := 0
	for _, line := range strings.Split(src, "\n") {
		len := len(line)
		if len > maxLen {
			maxLen = len
		}
	}
	fmt.Print(strings.Repeat("%", maxLen))
	fmt.Println(src)
	fmt.Println(strings.Repeat("%", maxLen))
	fmt.Println()
}

func printPath(path string) {
	fmt.Println(":: " + path)
}
