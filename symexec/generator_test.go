package symexec

import "testing"

func TestInitSmtFloat_Normal(t *testing.T) {
	want := "a_bits := uint64(0x41f4c0012080043d) // 5570040328.001035\na := math.Float64frombits(a_bits)"
	got, err := initSmtFloat64("a", "(fp #b0 #b10000011111 #x4c0012080043d)")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestInitSmtFloat_PosZero(t *testing.T) {
	want := "a := 0.0"
	got, err := initSmtFloat64("a", "(_ +zero 11 53)")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestInitSmtFloat_NegZero(t *testing.T) {
	want := "a := 0.0\na *= -1.0"
	got, err := initSmtFloat64("a", "(_ -zero 11 53)")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestInitSmtFloat_NaN(t *testing.T) {
	want := "a := math.NaN()"
	got, err := initSmtFloat64("a", "(_ NaN 11 53)")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestInitSmtFloat_PosInf(t *testing.T) {
	want := "a := math.Inf(1)"
	got, err := initSmtFloat64("a", "(_ +oo 11 53)")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestInitSmtFloat_NegInf(t *testing.T) {
	want := "a := math.Inf(-1)"
	got, err := initSmtFloat64("a", "(_ -oo 11 53)")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}
