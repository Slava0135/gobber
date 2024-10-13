package constraints

import (
	"strconv"

	"github.com/aclements/go-z3/z3"
)

func PushPopIncrementality() {
	src := `
func pushPopIncrementality(j int) int {
    result := j

    for i := 1; i <= 10; i++ {
        result += i
    }

    if result%2 == 0 {
        result++
    }
}`
	printSrc(src)

	ctx := z3.NewContext(nil)
	intSort := ctx.BVSort(intSize)
	j := ctx.BVConst("j", intSize)

	int0 := ctx.FromInt(0, intSort).(z3.BV)
	int1 := ctx.FromInt(1, intSort).(z3.BV)
	int2 := ctx.FromInt(2, intSort).(z3.BV)

	solver := z3.NewSolver(ctx)

	iInit := ctx.Const("i_1", intSort).(z3.BV)
	resultInit := ctx.Const("result_1", intSort).(z3.BV)
	solver.Assert(iInit.Eq(int1))
	solver.Assert(resultInit.Eq(j))
	resultPrev := resultInit
	iPrev := iInit
	for iter := 2; iter <= 10; iter++ {
		i := ctx.BVConst("i_"+strconv.Itoa(iter), intSize)
		result := ctx.BVConst("result_"+strconv.Itoa(iter), intSize)
		solver.Assert(i.Eq(iPrev.Add(int1)))
		solver.Assert(result.Eq(resultPrev.Add(i)))
		iPrev = i
		resultPrev = result
	}
	result := resultPrev

	solveIncrement(solver, "result%2 == 0", result.SToInt().Mod(int2.SToInt()).Eq(int0.SToInt()))
	solveIncrement(solver, "result%2 != 0", result.SToInt().Mod(int2.SToInt()).NE(int0.SToInt()))
}
