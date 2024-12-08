package main

import (
	"math"
)

type DoubleExamples struct{}

func (d *DoubleExamples) CompareSum(a, b float64) float64 {
	z := a + b
	if z > 5.6 {
		return 1.0
	} else {
		return 0.0
	}
}

func (d *DoubleExamples) Compare(a, b float64) float64 {
	if a > b {
		return 1.0
	} else {
		return 0.0
	}
}

func (d *DoubleExamples) CompareWithDiv(a, b float64) float64 {
	z := a + 0.5
	if (a / z) > b {
		return 1.0
	} else {
		return 0.0
	}
}

func (d *DoubleExamples) SimpleSum(a, b float64) float64 {
	if math.IsNaN(a + b) {
		return 0.0
	}
	c := a + 1.1
	if b+c > 10.1 && b+c < 11.125 {
		return 1.1
	} else {
		return 1.2
	}
}

func (d *DoubleExamples) Sum(a, b float64) float64 {
	if math.IsNaN(a + b) {
		return 0.0
	}
	c := a + 0.123124
	if b+c > 11.123124 && b+c < 11.125 {
		return 1.1
	} else {
		return 1.2
	}
}

func (d *DoubleExamples) SimpleMul(a, b float64) float64 {
	if math.IsNaN(a * b) {
		return 0.0
	}
	if a*b > 33.1 && a*b < 33.875 {
		return 1.1
	} else {
		return 1.2
	}
}

func (d *DoubleExamples) Mul(a, b float64) float64 {
	if math.IsNaN(a * b) {
		return 0.0
	}
	if a*b > 33.32 && a*b < 33.333 {
		return 1.1
	} else if a*b > 33.333 && a*b < 33.7592 {
		return 1.2
	} else {
		return 1.3
	}
}

func (d *DoubleExamples) CheckNonInteger(a float64) float64 {
	if a > 0.1 && a < 0.9 {
		return 1.0
	}
	return 0.0
}

func (d *DoubleExamples) Div(a, b, c float64) float64 {
	return (a + b) / c
}

func (d *DoubleExamples) SimpleEquation(a float64) int {
	if a+a+a-9 == a+3 {
		return 0
	} else {
		return 1
	}
}

func (d *DoubleExamples) SimpleNonLinearEquation(a float64) int {
	if 3*a-9 == a+3 {
		return 0
	} else {
		return 1
	}
}

func (d *DoubleExamples) CheckNaN(dValue float64) int {
	if dValue < 0 {
		return -1
	}
	if dValue > 0 {
		return 1
	}
	if dValue == 0 {
		return 0
	}
	// NaN
	return 100
}

func (d *DoubleExamples) UnaryMinus(dValue float64) int {
	if -dValue < 0 {
		return -1
	}
	return 0
}

func (d *DoubleExamples) DoubleInfinity(dValue float64) int {
	if dValue == math.Inf(1) {
		return 1
	}
	if dValue == math.Inf(-1) {
		return 2
	}
	return 3
}
