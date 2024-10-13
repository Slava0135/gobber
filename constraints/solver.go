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

func solveIncrementWithAssumptions(solver *z3.Solver, path string, assumptions []Assumption, asserts ...z3.Bool) {
	solver.Push()
	defer solver.Pop()

	printPath(path)
	for _, v := range asserts {
		solver.Assert(v)
	}

	assume := func(a []Assumption) {
		fmt.Print(":: assume ::")
		for _, v := range a {
			fmt.Print(" " + v.ref.String())
			solver.AssertAndTrack(v.assume, v.ref)
		}
		fmt.Println()
	}

	solver.Push()
	assume(assumptions)
	sat, err := solver.Check()
	if err != nil {
		panic(err)
	}
	if sat {
		fmt.Println(solver.Model())
		solver.Pop()
	} else {
		unsatCore := solver.GetUnsatCore()
		solver.Pop()
		var remaining []Assumption
		remaining = append(remaining, assumptions...)
		for !sat {
			fmt.Println("unsat core:", unsatCore)
			var nextRemaining []Assumption
			for i := range remaining {
				if remaining[i].ref.String() != unsatCore[len(unsatCore)-1].String() {
					nextRemaining = append(nextRemaining, remaining[i])
				}
			}
			remaining = nextRemaining

			solver.Push()
			assume(remaining)
			sat, err = solver.Check()
			if err != nil {
				panic(err)
			}
			if sat {
				fmt.Println(solver.Model())
				solver.Pop()
				return
			}
			unsatCore = solver.GetUnsatCore()
			solver.Pop()
		}
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
