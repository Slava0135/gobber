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

func TestStatic_Arrays(t *testing.T) {
	check(t, []string{}, AnalyzeFileStatic("arrays.go"))
}

func TestStatic_Complex(t *testing.T) {
	check(t, []string{"complexComparison"}, AnalyzeFileStatic("complex.go"))
}

func TestStatic_Numbers(t *testing.T) {
	check(t, []string{}, AnalyzeFileStatic("numbers.go"))
}

func TestStatic_PushPop(t *testing.T) {
	check(t, []string{"pushPopIncrementality"}, AnalyzeFileStatic("push_pop.go"))
}

func TestStatic_SoftConstraints(t *testing.T) {
	check(t, []string{}, AnalyzeFileStatic("softconstraints.go"))
}

func TestDynamic_Arrays(t *testing.T) {
	check(t, []string{}, AnalyzeFileDynamic("arrays.go"))
}

func TestDynamic_Complex(t *testing.T) {
	check(t, []string{"complexComparison"}, AnalyzeFileDynamic("complex.go"))
}

func TestDynamic_Numbers(t *testing.T) {
	check(t, []string{}, AnalyzeFileDynamic("numbers.go"))
}

func TestDynamic_PushPop(t *testing.T) {
	check(t, []string{}, AnalyzeFileDynamic("push_pop.go"))
}

func TestDynamic_SoftConstraints(t *testing.T) {
	check(t, []string{}, AnalyzeFileDynamic("softconstraints.go"))
}
