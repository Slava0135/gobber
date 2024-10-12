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
    floatSort := ctx.FloatSort(8, 24)
    x, ok := ctx.Const("x", floatSort).(z3.Float)
    if !ok {
        panic("not a float")
    }
    y, ok := ctx.Const("y", floatSort).(z3.Float)
    if !ok {
        panic("not a float")
    }

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
