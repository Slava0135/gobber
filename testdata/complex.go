package main

func basicComplexOperations(a complex128, b complex128) complex128 {
	if real(a) > real(b) {
		return a + b
	} else if imag(a) > imag(b) {
		return a - b
	}
	return a * b
}

func complexMagnitude(a complex128) float64 {
	magnitude := real(a)*real(a) + imag(a)*imag(a)
	return magnitude
}

func complexComparison(a complex128, b complex128) string {
	magA := complexMagnitude(a)
	magB := complexMagnitude(b)

	if magA > magB {
		return "Magnitude of a is greater than b"
	} else if magA < magB {
		return "Magnitude of b is greater than a"
	}
	return "Magnitudes are equal"
}

func complexOperations(a complex128, b complex128) complex128 {
	if real(a) == 0 && imag(a) == 0 {
		return b
	} else if real(b) == 0 && imag(b) == 0 {
		return a
	} else if real(a) > real(b) {
		return a / b
	}
	return a + b
}

func nestedComplexOperations(a complex128, b complex128) complex128 {
    if real(a) < 0 {
        if imag(a) < 0 {
            return a * b
        }
        return a + b
    }

    if imag(b) < 0 {
        return a - b
    }
    return a + b
}