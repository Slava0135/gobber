package symexec

import "testing"

func TestInitSmtFloat_Normal(t *testing.T) {
	want := "a := 5570040328.001035"
	got, err := initSmtFloat64("a", "fp #b0 #b10000011111 #x4c0012080043d")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}
