package constraints

import "github.com/aclements/go-z3/z3"

type Complex struct {
	real z3.Float
	imag z3.Float
}

func complexConst(ctx *z3.Context, name string, floatSort z3.Sort) Complex {
	return Complex{
		real: ctx.Const(name + ".REAL", floatSort).(z3.Float),
		imag: ctx.Const(name + ".IMAG", floatSort).(z3.Float),
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
