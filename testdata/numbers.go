package main

func integerOperations(a int, b int) int {
	if a > b {
		return a + b
	} else if a < b {
		return a - b
	} else {
		return a * b
	}
}

func floatOperations(x float64, y float64) float64 {
	if x > y {
		return x / y
	} else if x < y {
		return x * y
	}
	return 0.0
}

func mixedOperations(a int, b float64) float64 {
	var result float64

	if a%2 == 0 {
		result = float64(a) + b
	} else {
		result = float64(a) - b
	}

	if result < 10 {
		result *= 2
	} else {
		result /= 2
	}

	return result
}

func nestedConditions(a int, b float64) float64 {
	if a < 0 {
		if b < 0 {
			return float64(a*-1) + b
		}
		return float64(a*-1) - b
	}
	return float64(a) + b
}

func bitwiseOperations(a int, b int) int {
	if a&1 == 0 && b&1 == 0 {
		return a | b
	} else if a&1 == 1 && b&1 == 1 {
		return a & b
	}
	return a ^ b
}

func advancedBitwise(a int, b int) int {
	if a > b {
		return a << 1
	} else if a < b {
		return b >> 1
	}
	return a ^ b
}

func combinedBitwise(a int, b int) int {
	if a&b == 0 {
		return a | b
	} else {
		result := a & b
		if result > 10 {
			return result ^ b
		}
		return result
	}
}

func nestedBitwise(a int, b int) int {
	if a < 0 {
		return -1
	}

	if b < 0 {
		return a ^ 0
	}

	if a&b == 0 {
		return a | b
	} else {
		return a & b
	}
}