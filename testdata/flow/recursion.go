package main

import (
	"errors"
)

type MathExamples struct{}

func (m *MathExamples) Factorial(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("IllegalArgumentException")
	}
	if n == 0 {
		return 1, nil
	}
	result, _ := m.Factorial(n - 1)
	return n * result, nil
}

func (m *MathExamples) Fib(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("IllegalArgumentException")
	}
	if n == 0 {
		return 0, nil
	}
	if n == 1 {
		return 1, nil
	}
	fib1, _ := m.Fib(n - 1)
	fib2, _ := m.Fib(n - 2)
	return fib1 + fib2, nil
}

func (m *MathExamples) Sum(fst, snd int) int {
	if snd == 0 {
		return fst
	}
	sign := int(float64(snd) / float64(abs(snd)))
	return m.Sum(fst+sign, snd-sign)
}

func (m *MathExamples) Pow(a, n int) (int, error) {
	if n < 0 {
		return 0, errors.New("IllegalArgumentException")
	}
	if n == 0 {
		return 1, nil
	}
	if n%2 == 1 {
		result, _ := m.Pow(a, n-1)
		return result * a, nil
	} else {
		b, _ := m.Pow(a, n/2)
		return b * b, nil
	}
}

func (m *MathExamples) InfiniteRecursion(i int) {
	m.InfiniteRecursion(i + 1)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
