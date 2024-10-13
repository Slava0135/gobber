package constraints

import (
	"fmt"

	"github.com/aclements/go-z3/z3"
)

func CompareAndIncrement() {
	src := `
func compareAndIncrement(a, b int) int {
    if a > b {
        c := a + 1

        if (c > b) {
            return 1
        } else {
            return -1
        }
    }

    return 42
}`
	printSrc(src)

	ctx := z3.NewContext(nil)
	intSort := ctx.BVSort(intSize)
	a := ctx.BVConst("a", intSize)
	b := ctx.BVConst("b", intSize)
	c := ctx.BVConst("c", intSize)

	int0 := ctx.FromInt(0, intSort).(z3.BV)
	int1 := ctx.FromInt(1, intSort).(z3.BV)
	int2 := ctx.FromInt(2, intSort).(z3.BV)
	int10 := ctx.FromInt(10, intSort).(z3.BV)

	solver := z3.NewSolver(ctx)
	assumptions := []Assumption{
		{a.SToInt().GE(int0.SToInt()), ctx.BoolConst("a >= 0")},
		{b.SToInt().GE(int0.SToInt()), ctx.BoolConst("b >= 0")},
		{c.SToInt().GE(int0.SToInt()), ctx.BoolConst("c >= 0")},
		{a.SToInt().LT(int2.SToInt()), ctx.BoolConst("a < 2")},
		{b.SToInt().LT(int2.SToInt()), ctx.BoolConst("b < 2")},
		{c.SToInt().LT(int2.SToInt()), ctx.BoolConst("c < 2")},
		{a.SToInt().LT(int10.SToInt()), ctx.BoolConst("a < 10")},
		{b.SToInt().LT(int10.SToInt()), ctx.BoolConst("b < 10")},
		{c.SToInt().LT(int10.SToInt()), ctx.BoolConst("c < 10")},
	}

	{
		solver.Assert(a.SToInt().GT(b.SToInt()))
		solver.Assert(c.Eq(a.Add(int1)))
		solver.Assert(c.SToInt().GT(b.SToInt()))
		unsatCore := solveIncrementWithAssumptions(solver, "(a > b) && (c > b)", assumptions)
		var remaining []Assumption
		remaining = append(remaining, assumptions...)
		for unsatCore != nil {
			fmt.Println("unsat core:", unsatCore)
			var nextRemaining []Assumption 
			for i := range remaining {
				if assumptions[i].ref.String() != unsatCore[len(unsatCore)-1].String() {
					nextRemaining = append(nextRemaining, remaining[i])
				}
			}
			remaining = nextRemaining
			unsatCore = solveIncrementWithAssumptions(solver, "(a > b) && (c > b)", remaining)
		}
		fmt.Println(solver.Model())
	}

	solver.Reset()
	solveIncrementWithAssumptions(solver, "(a > b) && (c <= b)", assumptions,
		a.SToInt().GT(b.SToInt()),
		c.Eq(a.Add(int1)),
		c.SToInt().LE(b.SToInt()),
	)

	solver.Reset()
	solveIncrementWithAssumptions(solver, "a <= b", assumptions, a.SToInt().GE(b.SToInt()))
}
