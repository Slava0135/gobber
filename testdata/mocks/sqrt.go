package main

import (
	"math"
	"slava0135/gobber/symbolic"
)

func MockSqrt(f float64) float64 {
	s := math.Sqrt(f)
	if s > 4.0 {
		return 1.0
	}
	return 0.0
}

func math__Sqrt(x float64) float64 {
	res := symbolic.MakeSymbolic[float64]()
	symbolic.Assume(res >= 0)
	symbolic.Assume(math.Abs(res*res-x) < 1e-6)
	return res
}
