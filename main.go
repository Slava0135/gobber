package main

import (
	"slava0135/gobber/constraints"
	"slava0135/gobber/subtypes"
)

func main() {
	constraints.IntegerOperations()
	constraints.FloatOperations()
	constraints.MixedOperations()
	constraints.NestedConditions()
	constraints.BitwiseOperations()
	constraints.AdvancedBitwise()
	constraints.CombinedBitwise()
	constraints.NestedBitwise()

	constraints.BasicComplexOperations()
	constraints.ComplexMagnitude()
	constraints.ComplexComparison()
	constraints.ComplexOperations()
	constraints.NestedComplexOperations()

	constraints.CompareElement()
	constraints.CompareAge()

	constraints.PushPopIncrementality()

	constraints.CompareAndIncrement()

	subtypes.SubclassesExample()
	subtypes.SubtypesExample()
	subtypes.NaiveTypeSolver()
}
