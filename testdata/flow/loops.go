package main

type Loops struct{}

func (l *Loops) LoopWithConcreteBound(n int) int {
	result := 0
	for i := 0; i < 10; i++ {
		result += i
	}
	return result
}

func (l *Loops) LoopWithSymbolicBound(n int) int {
	if n > 10 {
		panic("Assumption violated: n should be less than or equal to 10")
	}

	result := 0
	for i := 0; i < n; i++ {
		result += i
	}
	return result
}

func (l *Loops) LoopWithConcreteBoundAndSymbolicBranching(condition bool) int {
	result := 0
	for i := 0; i < 10; i++ {
		if condition && i%2 == 0 {
			result += i
		}
	}
	return result
}

func (l *Loops) LoopWithSymbolicBoundAndSymbolicBranching(n int, condition bool) int {
	if n > 10 {
		panic("Assumption violated: n should be less than or equal to 10")
	}

	result := 0
	for i := 0; i < n; i++ {
		if condition && i%2 == 0 {
			result += i
		}
	}
	return result
}

func (l *Loops) LoopWithSymbolicBoundAndComplexControlFlow(n int, condition bool) int {
	if n > 10 {
		panic("Assumption violated: n should be less than or equal to 10")
	}

	result := 0
	for i := 0; i < n; i++ {
		if condition && i == 3 {
			break
		}
		if i%2 != 0 {
			continue
		}
		result += i
	}
	return result
}

func (l *Loops) ForCycle(x int) int {
	for i := 0; i < x; i++ {
		if i > 5 {
			return 1
		}
	}
	return -1
}

func (l *Loops) ForCycleFour(x int) int {
	for i := 0; i < x; i++ {
		if i > 4 {
			return 1
		}
	}
	return -1
}

func (l *Loops) FiniteCycle(x int) int {
	for {
		if x%519 == 0 {
			break
		} else {
			x++
		}
	}
	return x
}

func (l *Loops) ForCycleFromJayHorn(x int) int {
	r := 0
	for i := 0; i < x; i++ {
		r += 2
	}
	return r
}

func (l *Loops) DivideByZeroCheckWithCycles(n int, x int) int {
	if n < 5 {
		panic("IllegalArgumentException: n < 5")
	}
	j := 0
	for i := 0; i < n; i++ {
		j += i
	}
	if x == 0 {
		panic("Division by zero")
	}
	j /= x
	for i := 0; i < n; i++ {
		j += i
	}
	return 1
}

func (l *Loops) MoveToException(x int) {
	if x < 400 {
		for i := x; i < 400; i++ {
			x++
		}
	}

	if x > 400 {
		for i := x; i > 400; i-- {
			x--
		}
	}

	if x == 400 {
		panic("IllegalArgumentException")
	}
}

func (l *Loops) WhileCycle(x int) int {
	i := 0
	sum := 0
	for i < x {
		sum += i
		i++
	}
	return sum
}

func (l *Loops) CallInnerWhile(value int) int {
	return l.InnerWhile(value, 42)
}

func (l *Loops) InnerLoop(value int) int {
	cycleDependedCondition := &CycleDependedCondition{}
	return cycleDependedCondition.TwoCondition(value)
}

func (l *Loops) InnerWhile(a int, border int) int {
	res := a
	for res >= border {
		res -= border
	}
	return res
}

func (l *Loops) LoopInsideLoop(x int) int {
	for i := x - 5; i < x; i++ {
		if i < 0 {
			return 2
		} else {
			for j := i; j < x+i; j++ {
				if j == 7 {
					return 1
				}
			}
		}
	}
	return -1
}

func (l *Loops) StructureLoop(x int) int {
	for i := 0; i < x; i++ {
		if i == 2 {
			return 1
		}
	}
	return -1
}

type CycleDependedCondition struct{}

func (c *CycleDependedCondition) TwoCondition(x int) int {
	for i := 0; i < x; i++ {
		if i > 2 && x == 4 {
			return 1
		}
	}
	return 0
}
