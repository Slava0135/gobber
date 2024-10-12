package constraints

import "github.com/aclements/go-z3/z3"

func CompareElement() {
	src := `
func compareElement(array []int, index int, value int) int {
    if index < 0 || index >= len(array) {
        return -1 // Индекс вне границ
    }
    element := array[index]
    if element > value {
        return 1 // Элемент больше
    } else if element < value {
        return -1 // Элемент меньше
    }
    return 0 // Элемент равен
}`
	printSrc(src)

	ctx := z3.NewContext(nil)

	int0 := ctx.FromInt(0, ctx.IntSort()).(z3.Int)

	array := ctx.ConstArray(ctx.IntSort(), int0)
	arrayLen := ctx.IntConst("arrayLen")
	index := ctx.IntConst("index")
	value := ctx.IntConst("value")
	element := ctx.IntConst("element")

	solver := z3.NewSolver(ctx)

	assertArrayLen := arrayLen.GE(int0)

	solve(solver, "(index < 0) || (index >= len(array))", assertArrayLen, index.LT(int0).Or(index.GE(arrayLen)))
	solve(solver, "(index >= 0) && (index < len(array)) && (element > value)",
		assertArrayLen,
		index.GE(int0),
		index.LT(arrayLen),
		element.Eq(array.Select(index).(z3.Int)),
		element.GT(value),
	)
	solve(solver, "(index >= 0) && (index < len(array)) && !(element > value) && (element < value)",
		assertArrayLen,
		index.GE(int0),
		index.LT(arrayLen),
		element.Eq(array.Select(index).(z3.Int)),
		element.GT(value).Not(),
		element.LT(value),
	)
	solve(solver, "(index >= 0) && (index < len(array)) && !(element > value) && !(element < value)",
		assertArrayLen,
		index.GE(int0),
		index.LT(arrayLen),
		element.Eq(array.Select(index).(z3.Int)),
		element.GT(value).Not(),
		element.LT(value).Not(),
	)
}
