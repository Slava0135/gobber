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
	checkDynamic(t, []string{"complexComparison"}, "complex.go")
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
