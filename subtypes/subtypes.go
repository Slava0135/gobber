package subtypes

import (
	"fmt"

	"github.com/aclements/go-z3/z3"
)

type TypeLatticeElement struct {
	uppers    []*TypeLatticeElement
	canCastTo *TypeLatticeElement
}

// Scala types
var (
	Any TypeLatticeElement = TypeLatticeElement{}

	AnyVal TypeLatticeElement = TypeLatticeElement{}
	Short  TypeLatticeElement = TypeLatticeElement{}
	Int    TypeLatticeElement = TypeLatticeElement{}
	Long   TypeLatticeElement = TypeLatticeElement{}
	Char   TypeLatticeElement = TypeLatticeElement{}

	AnyRef TypeLatticeElement = TypeLatticeElement{}
	String TypeLatticeElement = TypeLatticeElement{}
	Seq    TypeLatticeElement = TypeLatticeElement{}
	List   TypeLatticeElement = TypeLatticeElement{}
	Null   TypeLatticeElement = TypeLatticeElement{}

	All []*TypeLatticeElement = []*TypeLatticeElement{&Any, &AnyVal, &Short, &Int, &Long, &Char, &AnyRef, &String, &Seq, &List, &Null}
)

func (e *TypeLatticeElement) isSubclassOf(other *TypeLatticeElement) bool {
	if e == other {
		return true
	}
	var queue []*TypeLatticeElement
	queue = append(queue, e.uppers...)
	for len(queue) > 0 {
		next := queue[len(queue)-1]
		if next == other {
			return true
		}
		queue = queue[:len(queue)-1]
	}
	return false
}

func (e *TypeLatticeElement) isSubtypeOf(other *TypeLatticeElement) bool {
	if e.isSubclassOf(other) {
		return true
	}
	if e.canCastTo == nil {
		return false
	}
	return e.canCastTo.isSubtypeOf(other)
}

func init() {
	AnyVal.uppers = append(AnyVal.uppers, &Any)
	AnyRef.uppers = append(AnyVal.uppers, &Any)

	Short.uppers = append(Short.uppers, &AnyVal)
	Int.uppers = append(Int.uppers, &AnyVal)
	Long.uppers = append(Long.uppers, &AnyVal)
	Char.uppers = append(Char.uppers, &AnyVal)

	Short.canCastTo = &Int
	Char.canCastTo = &Int
	Int.canCastTo = &Long

	String.uppers = append(String.uppers, &AnyRef)
	Seq.uppers = append(Seq.uppers, &AnyRef)
	List.uppers = append(List.uppers, &Seq)
	Null.uppers = append(Null.uppers, &String, &Seq, &List)
}

func SubclassesExample() {
	fmt.Println("::", "[Int]   isSubclassOf  [Int]     >", Int.isSubclassOf(&Int))
	fmt.Println("::", "[Int]   isSubclassOf  [Long]    >", Int.isSubclassOf(&Long))
	fmt.Println("::", "[Long]  isSubclassOf  [Int]     >", Long.isSubclassOf(&Int))
	fmt.Println("::", "[Int]   isSubclassOf  [AnyVal]  >", Int.isSubclassOf(&AnyVal))
	fmt.Println("::", "[Int]   isSubclassOf  [Seq]     >", Int.isSubclassOf(&Seq))
	fmt.Println("::", "[List]  isSubclassOf  [Seq]     >", List.isSubclassOf(&Seq))
	fmt.Println()
}

func SubtypesExample() {
	fmt.Println("::", "[Int]   <:  [Int]     >", Int.isSubtypeOf(&Int))
	fmt.Println("::", "[Int]   <:  [Long]    >", Int.isSubtypeOf(&Long))
	fmt.Println("::", "[Long]  <:  [Int]     >", Long.isSubtypeOf(&Int))
	fmt.Println("::", "[Int]   <:  [AnyVal]  >", Int.isSubtypeOf(&AnyVal))
	fmt.Println("::", "[Int]   <:  [Seq]     >", Int.isSubtypeOf(&Seq))
	fmt.Println("::", "[List]  <:  [Seq]     >", List.isSubtypeOf(&Seq))
	fmt.Println()
}

func NaiveTypeSolver() {
	ctx := z3.NewContext(nil)

	typeSort := ctx.UninterpretedSort("Type")

	anyT := ctx.Const("Any", typeSort)

	anyValT := ctx.Const("AnyVal", typeSort)
	shortT := ctx.Const("Short", typeSort)
	intT := ctx.Const("Int", typeSort)
	longT := ctx.Const("Long", typeSort)
	charT := ctx.Const("Char", typeSort)

	anyRefT := ctx.Const("AnyRef", typeSort)
	stringT := ctx.Const("String", typeSort)
	seqT := ctx.Const("Seq", typeSort)
	listT := ctx.Const("List", typeSort)
	nullT := ctx.Const("Null", typeSort)

	isSubclassOf := ctx.FuncDecl("isSubclassOf", []z3.Sort{typeSort, typeSort}, ctx.BoolSort())

	solver := z3.NewSolver(ctx)

	allT := []z3.Value{anyT, anyValT, shortT, intT, longT, charT, anyRefT, stringT, seqT, listT, nullT}
	for i := range allT {
		for j := range allT {
			if All[i].isSubclassOf(All[j]) {
				solver.Assert(isSubclassOf.Apply(allT[i], allT[j]).(z3.Bool))
			} else {
				solver.Assert(isSubclassOf.Apply(allT[i], allT[j]).(z3.Bool).Not())
			}
		}
	}

	sat, err := solver.Check()
	if err != nil {
		panic(err)
	}
	if !sat {
		panic("unexpected unsat")
	}
	fmt.Println(solver.Model())

	solver.Push()
	fmt.Print(":: ", "[List]  isSubclassOf  [Seq]")
	solver.Assert(isSubclassOf.Apply(listT, seqT).(z3.Bool))
	sat, err = solver.Check()
	if err != nil {
		panic(err)
	}
	if !sat {
		panic("unexpected unsat")
	}
	fmt.Println("  > true")
	solver.Pop()

	solver.Push()
	fmt.Print(":: ", "[Long]  isSubclassOf  [Int]")
	solver.Assert(isSubclassOf.Apply(longT, intT).(z3.Bool))
	sat, err = solver.Check()
	if err != nil {
		panic(err)
	}
	if sat {
		panic("unexpected sat")
	}
	fmt.Println("  > false")
	solver.Pop()
}
