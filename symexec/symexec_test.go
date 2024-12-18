package symexec

import (
	"os"
	"slices"
	"testing"
)

func TestMain(m *testing.M) {
	os.Chdir("../testdata")
	m.Run()
}

func checkStatic(t *testing.T, shouldFail []string, filename string) {
	r := AnalyzeFileStatic(filename)
	for fn, ok := range r {
		if ok && slices.Contains(shouldFail, fn) {
			t.Errorf("'%s' should fail", fn)
		}
		if !ok && !slices.Contains(shouldFail, fn) {
			t.Errorf("'%s' should succeed", fn)
		}
	}
}

func checkDynamic(t *testing.T, shouldFail []string, filename string) {
	testcases := AnalyzeFileDynamic(filename)
	for fn, tc := range testcases {
		if tc != nil && slices.Contains(shouldFail, functionName(fn)) {
			t.Errorf("'%s' should fail", fn)
		}
		if tc == nil && !slices.Contains(shouldFail, functionName(fn)) {
			t.Errorf("'%s' should succeed", fn)
		}
	}
	GenerateTests(filename, testcases)
}

func TestStatic_Arrays(t *testing.T) {
	checkStatic(t, []string{}, "arrays.go")
}

func TestStatic_Complex(t *testing.T) {
	checkStatic(t, []string{"complexComparison"}, "complex.go")
}

func TestStatic_Numbers(t *testing.T) {
	checkStatic(t, []string{}, "numbers.go")
}

func TestStatic_PushPop(t *testing.T) {
	checkStatic(t, []string{"pushPopIncrementality"}, "push_pop.go")
}

func TestStatic_SoftConstraints(t *testing.T) {
	checkStatic(t, []string{}, "softconstraints.go")
}

func TestDynamic_Arrays(t *testing.T) {
	checkDynamic(t, []string{}, "arrays.go")
}

func TestDynamic_Complex(t *testing.T) {
	checkDynamic(t, []string{}, "complex.go")
}

func TestDynamic_Numbers(t *testing.T) {
	checkDynamic(t, []string{}, "numbers.go")
}

func TestDynamic_PushPop(t *testing.T) {
	checkDynamic(t, []string{}, "push_pop.go")
}

func TestDynamic_SoftConstraints(t *testing.T) {
	checkDynamic(t, []string{}, "softconstraints.go")
}

func TestDynamic_Primitives_Doubles(t *testing.T) {
	checkDynamic(t, []string{}, "primitives/doubles.go")
}

func TestDynamic_Primitives_Overflow(t *testing.T) {
	checkDynamic(t, []string{}, "primitives/overflow.go")
}

func TestDynamic_Operators_Bit(t *testing.T) {
	checkDynamic(t, []string{}, "operators/bit.go")
}

func TestDynamic_Objects_RecursiveStruct(t *testing.T) {
	checkDynamic(t, []string{}, "objects/recursiveStruct.go")
}

func TestDynamic_Objects_WithPrimitives(t *testing.T) {
	checkDynamic(t, []string{}, "objects/withPrimitives.go")
}

func TestDynamic_Objects_WithReference(t *testing.T) {
	checkDynamic(t, []string{}, "objects/withReference.go")
}

func TestDynamic_Invokes_SimpleCalls(t *testing.T) {
	checkDynamic(t, []string{}, "invokes/simpleCalls.go")
}

func TestDynamic_Flow_Loops(t *testing.T) {
	checkDynamic(t, []string{}, "flow/loops.go")
}

func TestDynamic_Flow_Recursion(t *testing.T) {
	checkDynamic(t, []string{}, "flow/recursion.go")
}

func TestDynamic_Arrays_ArrayOfArrays(t *testing.T) {
	checkDynamic(t, []string{}, "arrays/arrayOfArrays.go")
}

func TestDynamic_Arrays_ArrayOfObjects(t *testing.T) {
	checkDynamic(t, []string{}, "arrays/arrayOfObjects.go")
}

func TestDynamic_Arrays_ArrayOverwriteValue(t *testing.T) {
	checkDynamic(t, []string{}, "arrays/arrayOverwriteValue.go")
}

func TestDynamic_Arrays_PrimitiveArrays(t *testing.T) {
	checkDynamic(t, []string{}, "arrays/primitiveArrays.go")
}

func TestDynamic_Mocks_Sqrt(t *testing.T) {
	checkDynamic(t, []string{}, "mocks/sqrt.go")
}
