package constraints

import (
	"fmt"

	"github.com/aclements/go-z3/z3"
)

func IntegerOperations() {
	ctx := z3.NewContext(nil)
	a := ctx.IntConst("a")
	b := ctx.IntConst("b")

	solver := z3.NewSolver(ctx)

	fmt.Println(":: a > b")
	solver.Assert(a.GT(b))
	sat, err := solver.Check()
	if err != nil {
		panic(err)
	}
	if !sat {
		panic("unexpected unsat")
	}
	fmt.Println(solver.Model())
	solver.Reset()

	fmt.Println(":: !(a > b) && (a < b)")
	solver.Assert(a.GT(b).Not().And(a.LT(b)))
	sat, err = solver.Check()
	if err != nil {
		panic(err)
	}
	if !sat {
		panic("unexpected unsat")
	}
	fmt.Println(solver.Model())
	solver.Reset()

	fmt.Println(":: !(a > b) && !(a < b)")
	solver.Assert(a.GT(b).Not().And(a.LT(b).Not()))
	sat, err = solver.Check()
	if err != nil {
		panic(err)
	}
	if !sat {
		panic("unexpected unsat")
	}
	fmt.Println(solver.Model())
	solver.Reset()
}

func FloatOperations() {
	ctx := z3.NewContext(nil)
	floatSort := ctx.FloatSort(11, 53)
	x := ctx.Const("x", floatSort).(z3.Float)
	y := ctx.Const("y", floatSort).(z3.Float)

	solver := z3.NewSolver(ctx)

	fmt.Println(":: x > y")
	solver.Assert(x.GT(y))
	sat, err := solver.Check()
	if err != nil {
		panic(err)
	}
	if !sat {
		panic("unexpected unsat")
	}
	fmt.Println(solver.Model())
	solver.Reset()

	fmt.Println(":: !(x > y) && (x < y)")
	solver.Assert(x.GT(y).Not().And(x.LT(y)))
	sat, err = solver.Check()
	if err != nil {
		panic(err)
	}
	if !sat {
		panic("unexpected unsat")
	}
	fmt.Println(solver.Model())
	solver.Reset()

	fmt.Println(":: !(x > y) && !(x < y)")
	solver.Assert(x.GT(y).Not().And(x.LT(y).Not()))
	sat, err = solver.Check()
	if err != nil {
		panic(err)
	}
	if !sat {
		panic("unexpected unsat")
	}
	fmt.Println(solver.Model())
	solver.Reset()
}

func MixedOperations() {
	ctx := z3.NewContext(nil)
	floatSort := ctx.FloatSort(11, 53)
	a := ctx.BVConst("a", 64)
	b := ctx.Const("b", floatSort).(z3.Float)
	result := ctx.Const("result", floatSort).(z3.Float)

	solver := z3.NewSolver(ctx)

	int0 := ctx.FromInt(0, ctx.IntSort()).(z3.Int)
	int2 := ctx.FromInt(2, ctx.IntSort()).(z3.Int)
	float10 := ctx.FromInt(10, floatSort).(z3.Float)

	fmt.Println(":: (a % 2 == 0) && (result < 10)")
	solver.Assert(a.SToInt().Mod(int2).Eq(int0))
	solver.Assert(result.Eq(a.SToFloat(floatSort).Add(b)))
	solver.Assert(result.LT(float10))
	sat, err := solver.Check()
	if err != nil {
		panic(err)
	}
	if !sat {
		panic("unexpected unsat")
	}
	fmt.Println(solver.Model())
	solver.Reset()

	fmt.Println(":: (a % 2 /= 0) && (result < 10)")
	solver.Assert(a.SToInt().Mod(int2).NE(int0))
	solver.Assert(result.Eq(a.SToFloat(floatSort).Add(b)))
	solver.Assert(result.LT(float10))
	sat, err = solver.Check()
	if err != nil {
		panic(err)
	}
	if !sat {
		panic("unexpected unsat")
	}
	fmt.Println(solver.Model())
	solver.Reset()

	fmt.Println(":: (a % 2 == 0) && (result >= 10)")
	solver.Assert(a.SToInt().Mod(int2).Eq(int0))
	solver.Assert(result.Eq(a.SToFloat(floatSort).Add(b)))
	solver.Assert(result.GE(float10))
	sat, err = solver.Check()
	if err != nil {
		panic(err)
	}
	if !sat {
		panic("unexpected unsat")
	}
	fmt.Println(solver.Model())
	solver.Reset()

	fmt.Println(":: (a % 2 /= 0) && (result >= 10)")
	solver.Assert(a.SToInt().Mod(int2).NE(int0))
	solver.Assert(result.Eq(a.SToFloat(floatSort).Add(b)))
	solver.Assert(result.GE(float10))
	sat, err = solver.Check()
	if err != nil {
		panic(err)
	}
	if !sat {
		panic("unexpected unsat")
	}
	fmt.Println(solver.Model())
	solver.Reset()
}
