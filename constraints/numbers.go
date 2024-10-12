package constraints

import (
	"fmt"

	"github.com/aclements/go-z3/z3"
)

const IntSize = 64

func IntegerOperations() {
	ctx := z3.NewContext(nil)
	a := ctx.IntConst("a")
	b := ctx.IntConst("b")

	solver := z3.NewSolver(ctx)

	solve(solver, "a > b", a.GT(b))
	solve(solver, "!(a > b) && (a < b)", a.GT(b).Not(), a.LT(b))
	solve(solver, "!(a > b) && !(a < b)", a.GT(b).Not(), a.LT(b).Not())
}

func FloatOperations() {
	ctx := z3.NewContext(nil)
	floatSort := ctx.FloatSort(11, 53)
	x := ctx.Const("x", floatSort).(z3.Float)
	y := ctx.Const("y", floatSort).(z3.Float)

	solver := z3.NewSolver(ctx)

	solve(solver, "x > y", x.GT(y))
	solve(solver, "!(x > y) && (x < y)", x.GT(y).Not(), x.LT(y))
	solve(solver, "!(x > y) && !(x < y)", x.GT(y).Not(), x.LT(y).Not())
}

func MixedOperations() {
	ctx := z3.NewContext(nil)
	floatSort := ctx.FloatSort(11, 53)
	a := ctx.BVConst("a", IntSize)
	b := ctx.Const("b", floatSort).(z3.Float)
	result := ctx.Const("result", floatSort).(z3.Float)

	int0 := ctx.FromInt(0, ctx.IntSort()).(z3.Int)
	int2 := ctx.FromInt(2, ctx.IntSort()).(z3.Int)
	float10 := ctx.FromInt(10, floatSort).(z3.Float)
    
	solver := z3.NewSolver(ctx)

	solve(solver, "(a % 2 == 0) && (result < 10)",
		a.SToInt().Mod(int2).Eq(int0),
		result.Eq(a.SToFloat(floatSort).Add(b)),
		result.LT(float10),
	)

	solve(solver, "(a % 2 == 0) && (result >= 10)",
		a.SToInt().Mod(int2).Eq(int0),
		result.Eq(a.SToFloat(floatSort).Add(b)),
		result.GE(float10),
	)

	solve(solver, "(a % 2 /= 0) && (result < 10)",
		a.SToInt().Mod(int2).Eq(int0),
		result.NE(a.SToFloat(floatSort).Add(b)),
		result.LT(float10),
	)

	solve(solver, "(a % 2 /= 0) && (result >= 10)",
		a.SToInt().Mod(int2).Eq(int0),
		result.NE(a.SToFloat(floatSort).Add(b)),
		result.GE(float10),
	)
}

func NestedConditions() {
    ctx := z3.NewContext(nil)
	floatSort := ctx.FloatSort(11, 53)
	a := ctx.BVConst("a", IntSize)
	b := ctx.Const("b", floatSort).(z3.Float)

    int0 := ctx.FromInt(0, ctx.IntSort()).(z3.Int)
	float0 := ctx.FromInt(0, floatSort).(z3.Float)

    solver := z3.NewSolver(ctx)

    solve(solver, "a < 0 && b < 0", a.SToInt().LT(int0), b.LT(float0))
    solve(solver, "a < 0 && b >= 0", a.SToInt().LT(int0), b.GE(float0))
    solve(solver, "a > 0", a.SToInt().GE(int0))
}

func solve(solver *z3.Solver, path string, asserts ...z3.Bool) {
	fmt.Println(":: " + path)
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
