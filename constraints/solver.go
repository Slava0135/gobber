package constraints

import (
	"fmt"
	"strings"

	"github.com/aclements/go-z3/z3"
)

const intSize = 64

type Assumption struct {
	assume z3.Bool
	ref    z3.Bool
}

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

func solveIncrementWithAssumptions(solver *z3.Solver, path string, assumptions []Assumption, asserts ...z3.Bool) []z3.Bool {
	solver.Push()
	defer solver.Pop()
	printPath(path)
	for _, v := range asserts {
		solver.Assert(v)
	}
	fmt.Print(":: assume ::")
	for _, v := range assumptions {
		fmt.Print(" " + v.ref.String())
		solver.AssertAndTrack(v.assume, v.ref)
	}
	fmt.Println()
	sat, err := solver.Check()
	if err != nil {
		panic(err)
	}
	if sat {
		fmt.Println(solver.Model())
		return nil
	} else {
		return solver.GetUnsatCore()
	}
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
