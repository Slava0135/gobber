package main

import (
	"math"
)

func MockSqrt(f float64) float64 {
	s := math.Sqrt(f)
	if s > 4.0 {
		return 1.0
	}
	return 0.0
}
