package constraints

import (
	"github.com/aclements/go-z3/z3"
)

func IntegerOperations() {
	src := `
func integerOperations(a int, b int) int {
	if a > b {
		return a + b
	} else if a < b {
		return a - b
	} else {
		return a * b
	}
}`
	printSrc(src)

	ctx := z3.NewContext(nil)
	a := ctx.IntConst("a")
	b := ctx.IntConst("b")

	solver := z3.NewSolver(ctx)

	solve(solver, "a > b", a.GT(b))
	solve(solver, "!(a > b) && (a < b)", a.GT(b).Not(), a.LT(b))
	solve(solver, "!(a > b) && !(a < b)", a.GT(b).Not(), a.LT(b).Not())
}

func FloatOperations() {
	src := `
func floatOperations(x float64, y float64) float64 {
	if x > y {
		return x / y
	} else if x < y {
		return x * y
	}
	return 0.0
}`
	printSrc(src)

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
	src := `
func mixedOperations(a int, b float64) float64 {
	var result float64

	if a%2 == 0 {
		result = float64(a) + b
	} else {
		result = float64(a) - b
	}

	if result < 10 {
		result *= 2
	} else {
		result /= 2
	}

	return result
}`
	printSrc(src)

	ctx := z3.NewContext(nil)
	floatSort := ctx.FloatSort(11, 53)
	a := ctx.BVConst("a", intSize)
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

	solve(solver, "(a % 2 != 0) && (result < 10)",
		a.SToInt().Mod(int2).Eq(int0),
		result.NE(a.SToFloat(floatSort).Add(b)),
		result.LT(float10),
	)

	solve(solver, "(a % 2 != 0) && (result >= 10)",
		a.SToInt().Mod(int2).Eq(int0),
		result.NE(a.SToFloat(floatSort).Add(b)),
		result.GE(float10),
	)
}

func NestedConditions() {
	src := `
func nestedConditions(a int, b float64) float64 {
	if a < 0 {
		if b < 0 {
			return float64(a*-1) + b
		}
		return float64(a*-1) - b
	}
	return float64(a) + b
}`
	printSrc(src)

	ctx := z3.NewContext(nil)
	floatSort := ctx.FloatSort(11, 53)
	a := ctx.BVConst("a", intSize)
	b := ctx.Const("b", floatSort).(z3.Float)

	int0 := ctx.FromInt(0, ctx.IntSort()).(z3.Int)
	float0 := ctx.FromInt(0, floatSort).(z3.Float)

	solver := z3.NewSolver(ctx)

	solve(solver, "a < 0 && b < 0", a.SToInt().LT(int0), b.LT(float0))
	solve(solver, "a < 0 && b >= 0", a.SToInt().LT(int0), b.GE(float0))
	solve(solver, "a > 0", a.SToInt().GE(int0))
}

func BitwiseOperations() {
	src := `
func bitwiseOperations(a int, b int) int {
	if a&1 == 0 && b&1 == 0 {
		return a | b
	} else if a&1 == 1 && b&1 == 1 {
		return a & b
	}
	return a ^ b
}`
	printSrc(src)

	ctx := z3.NewContext(nil)
	a := ctx.BVConst("a", intSize)
	b := ctx.BVConst("b", intSize)
	intSort := ctx.BVSort(intSize)

	int0 := ctx.FromInt(0, intSort).(z3.BV)
	int1 := ctx.FromInt(1, intSort).(z3.BV)

	solver := z3.NewSolver(ctx)

	solve(solver, "a&1 == 0 && b&1 == 0", a.And(int1).Eq(int0), b.And(int1).Eq(int0))
	solve(solver, "!(a&1 == 0 && b&1 == 0) && (a&1 == 1 && b&1 == 1)",
		a.And(int1).Eq(int0).And(b.And(int1).Eq(int0)).Not(),
		a.And(int1).Eq(int1).And(b.And(int1).Eq(int1)),
	)
	solve(solver, "!(a&1 == 0 && b&1 == 0) && !(a&1 == 1 && b&1 == 1)",
		a.And(int1).Eq(int0).And(b.And(int1).Eq(int0)).Not(),
		a.And(int1).Eq(int1).And(b.And(int1).Eq(int1)).Not(),
	)
}

func AdvancedBitwise() {
	src := `
func advancedBitwise(a int, b int) int {
	if a > b {
		return a << 1
	} else if a < b {
		return b >> 1
	}
	return a ^ b
}`
	printSrc(src)

	ctx := z3.NewContext(nil)
	a := ctx.BVConst("a", intSize)
	b := ctx.BVConst("b", intSize)

	solver := z3.NewSolver(ctx)

	solve(solver, "a > b", a.SToInt().GT(b.SToInt()))
	solve(solver, "!(a > b) && (a < b)", a.SToInt().GT(b.SToInt()).Not(), a.SToInt().LT(b.SToInt()))
	solve(solver, "!(a > b) && !(a < b)", a.SToInt().GT(b.SToInt()).Not(), a.SToInt().LT(b.SToInt()).Not())
}

func CombinedBitwise() {
	src := `
func combinedBitwise(a int, b int) int {
	if a&b == 0 {
		return a | b
	} else {
		result := a & b
		if result > 10 {
			return result ^ b
		}
		return result
	}
}`
	printSrc(src)

	ctx := z3.NewContext(nil)
	a := ctx.BVConst("a", intSize)
	b := ctx.BVConst("b", intSize)
	result := ctx.BVConst("result", intSize)
	intSort := ctx.BVSort(intSize)

	int0 := ctx.FromInt(0, intSort).(z3.BV)
	int10 := ctx.FromInt(10, intSort).(z3.BV)

	solver := z3.NewSolver(ctx)

	solve(solver, "a&b == 0", a.And(b).Eq(int0))
	solve(solver, "a&b != 0 && result > 10",
		a.And(b).NE(int0),
		result.Eq(a.And(b)),
		result.SToInt().GT(int10.SToInt()),
	)
	solve(solver, "a&b != 0 && result <= 10",
		a.And(b).NE(int0),
		result.Eq(a.And(b)),
		result.SToInt().LE(int10.SToInt()),
	)
}

func NestedBitwise() {
	src := `
func nestedBitwise(a int, b int) int {
	if a < 0 {
		return -1
	}

	if b < 0 {
		return a ^ 0
	}

	if a&b == 0 {
		return a | b
	} else {
		return a & b
	}
}`
	printSrc(src)

	ctx := z3.NewContext(nil)
	a := ctx.BVConst("a", intSize)
	b := ctx.BVConst("b", intSize)
	intSort := ctx.BVSort(intSize)

	int0 := ctx.FromInt(0, intSort).(z3.BV)

	solver := z3.NewSolver(ctx)

	solve(solver, "a < 0", a.SToInt().LT(int0.SToInt()))
	solve(solver, "(a >= 0) && (b < 0)", a.SToInt().GE(int0.SToInt()), b.SToInt().LT(int0.SToInt()))
	solve(solver, "(a >= 0) && (b >= 0) && (a&b == 0)",
		a.SToInt().GE(int0.SToInt()),
		b.SToInt().GE(int0.SToInt()),
		a.And(b).Eq(int0),
	)
	solve(solver, "(a >= 0) && (b >= 0) && (a&b != 0)",
		a.SToInt().GE(int0.SToInt()),
		b.SToInt().GE(int0.SToInt()),
		a.And(b).NE(int0),
	)
}
