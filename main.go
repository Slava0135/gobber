package main

import (
	"flag"

	"slava0135/gobber/constraints"
	"slava0135/gobber/subtypes"
	"slava0135/gobber/symexec"
)

func main() {
	runNumbers := flag.Bool("numbers", false, "solve numbers constraints")
	runComplex := flag.Bool("complex", false, "solve complex constraints")
	runArrays := flag.Bool("arrays", false, "solve arrays constraints")
	runPushPop := flag.Bool("pushpop", false, "solve constraints with push-pop incrementality")
	runSoft := flag.Bool("soft", false, "solve soft constraints")

	runSubtypes := flag.Bool("subtypes", false, "subtyping encoding")

	runStatic := flag.Bool("static", false, "static symbolic execution")

	flag.Parse()

	if *runNumbers {
		constraints.IntegerOperations()
		constraints.FloatOperations()
		constraints.MixedOperations()
		constraints.NestedConditions()
		constraints.BitwiseOperations()
		constraints.AdvancedBitwise()
		constraints.CombinedBitwise()
		constraints.NestedBitwise()
	}

	if *runComplex {
		constraints.BasicComplexOperations()
		constraints.ComplexMagnitude()
		constraints.ComplexComparison()
		constraints.ComplexOperations()
		constraints.NestedComplexOperations()
	}

	if *runArrays {
		constraints.CompareElement()
		constraints.CompareAge()
	}

	if *runPushPop {
		constraints.PushPopIncrementality()
	}

	if *runSoft {
		constraints.CompareAndIncrement()
	}

	if *runSubtypes {
		subtypes.SubclassesExample()
		subtypes.SubtypesExample()
		subtypes.NaiveTypeSolver()
	}

	if *runStatic {
		symexec.Static()
	}
}
