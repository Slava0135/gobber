package constraints

import "github.com/aclements/go-z3/z3"

type Complex struct {
	real z3.Float
	imag z3.Float
}

func complexConst(ctx *z3.Context, name string, floatSort z3.Sort) Complex {
	return Complex{
		real: ctx.Const(name+".REAL", floatSort).(z3.Float),
		imag: ctx.Const(name+".IMAG", floatSort).(z3.Float),
	}
}

func BasicComplexOperations() {
	src := `
func basicComplexOperations(a complex128, b complex128) complex128 {
	if real(a) > real(b) {
		return a + b
	} else if imag(a) > imag(b) {
		return a - b
	}
	return a * b
}`
	printSrc(src)

	ctx := z3.NewContext(nil)
	floatSort := ctx.FloatSort(11, 53)

	a := complexConst(ctx, "a", floatSort)
	b := complexConst(ctx, "b", floatSort)

	solver := z3.NewSolver(ctx)

	solve(solver, "real(a) > real(b)", a.real.GT(b.real))
	solve(solver, "(real(a) <= real(b)) && (imag(a) > imag(b))", a.real.GE(b.real), a.imag.LT(b.imag))
	solve(solver, "(real(a) <= real(b)) && (imag(a) <= imag(b))", a.real.GE(b.real), a.imag.GE(b.imag))
}

func ComplexMagnitude() {
	src := `
func complexMagnitude(a complex128) float64 {
	magnitude := real(a)*real(a) + imag(a)*imag(a)
	return magnitude
}`
	printSrc(src)
}

func ComplexOperations() {
	src := `
func complexOperations(a complex128, b complex128) complex128 {
	if real(a) == 0 && imag(a) == 0 {
		return b
	} else if real(b) == 0 && imag(b) == 0 {
		return a
	} else if real(a) > real(b) {
		return a / b
	}
	return a + b
}`
	printSrc(src)

	ctx := z3.NewContext(nil)
	floatSort := ctx.FloatSort(11, 53)

	a := complexConst(ctx, "a", floatSort)
	b := complexConst(ctx, "b", floatSort)

	float0 := ctx.FromInt(0, floatSort).(z3.Float)

	solver := z3.NewSolver(ctx)

	solve(solver, "real(a) == 0 && imag(a) == 0", a.real.Eq(float0), a.imag.Eq(float0))
	solve(solver, "!(real(a) == 0 && imag(a) == 0) && (real(b) == 0 && imag(b) == 0)",
		a.real.Eq(float0).And(a.imag.Eq(float0)).Not(),
		b.real.Eq(float0).And(b.imag.Eq(float0)),
	)
	solve(solver, "!(real(a) == 0 && imag(a) == 0) && !(real(b) == 0 && imag(b) == 0) && (real(a) > real(b))",
		a.real.Eq(float0).And(a.imag.Eq(float0)).Not(),
		b.real.Eq(float0).And(b.imag.Eq(float0)).Not(),
		a.real.GT(b.real),
	)
	solve(solver, "!(real(a) == 0 && imag(a) == 0) && !(real(b) == 0 && imag(b) == 0) && (real(a) <= real(b))",
		a.real.Eq(float0).And(a.imag.Eq(float0)).Not(),
		b.real.Eq(float0).And(b.imag.Eq(float0)).Not(),
		a.real.LE(b.real),
	)
}

func NestedComplexOperations() {
	src := `
func nestedComplexOperations(a complex128, b complex128) complex128 {
    if real(a) < 0 {
        if imag(a) < 0 {
            return a * b
        }
        return a + b
    }

    if imag(b) < 0 {
        return a - b
    }
    return a + b
}`
	printSrc(src)

	ctx := z3.NewContext(nil)
	floatSort := ctx.FloatSort(11, 53)

	a := complexConst(ctx, "a", floatSort)
	b := complexConst(ctx, "b", floatSort)

	float0 := ctx.FromInt(0, floatSort).(z3.Float)

	solver := z3.NewSolver(ctx)

	solve(solver, "(real(a) < 0) && (imag(a) < 0)", a.real.LT(float0), a.imag.LT(float0))
	solve(solver, "(real(a) < 0) && (imag(a) >= 0)", a.real.LT(float0), a.imag.GE(float0))
	solve(solver, "(real(a) >= 0) && (imag(b) < 0)", a.real.GE(float0), b.imag.LT(float0))
	solve(solver, "(real(a) >= 0) && (imag(b) >= 0)", a.real.GE(float0), b.imag.GE(float0))
}
