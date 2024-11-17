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

func check(t *testing.T, shouldFail []string, r map[string]bool) {
	for f, ok := range r {
		if ok && slices.Contains(shouldFail, f) {
			t.Errorf("'%s' should fail", f)
		}
		if !ok && !slices.Contains(shouldFail, f) {
			t.Errorf("'%s' should succeed", f)
		}
	}
}

func TestArrays(t *testing.T) {
	check(t, []string{}, AnalyzeFile("arrays.go"))
}

func TestComplex(t *testing.T) {
	check(t, []string{"complexComparison"}, AnalyzeFile("complex.go"))
}

func TestNumbers(t *testing.T) {
	check(t, []string{}, AnalyzeFile("numbers.go"))
}

func TestPushPop(t *testing.T) {
	check(t, []string{"pushPopIncrementality"}, AnalyzeFile("push_pop.go"))
}

func TestSoftConstraints(t *testing.T) {
	check(t, []string{}, AnalyzeFile("softconstraints.go"))
}
