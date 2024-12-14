package main

type PrimitiveArrays struct{}

func (p *PrimitiveArrays) DefaultIntValues() []int {
	array := make([]int, 3)
	if array[1] != 0 {
		return array
	}
	return array
}

func (p *PrimitiveArrays) DefaultDoubleValues() []float64 {
	array := make([]float64, 3)
	if array[1] != 0.0 {
		return array
	}
	return array
}

func (p *PrimitiveArrays) DefaultBooleanValues() []bool {
	array := make([]bool, 3)
	if !array[1] {
		return array
	}
	return array
}

func (p *PrimitiveArrays) ByteArray(a []byte, x byte) byte {
	if len(a) != 2 {
		return 255
	}
	a[0] = 5
	a[1] = x
	if a[0]+a[1] > 20 {
		return 1
	}
	return 0
}

func (p *PrimitiveArrays) ShortArray(a []int16, x int16) byte {
	if len(a) != 2 {
		return 255
	}
	a[0] = 5
	a[1] = x
	if a[0]+a[1] > 20 {
		return 1
	}
	return 0
}

func (p *PrimitiveArrays) CharArray(a []rune, x rune) byte {
	if len(a) != 2 {
		return 255
	}
	a[0] = 5
	a[1] = x

	sum := int(a[0]) + int(a[1])
	if sum > 20 {
		return 1
	}
	return 0
}

func (p *PrimitiveArrays) IntArray(a []int, x int) byte {
	if len(a) != 2 {
		return 255
	}
	a[0] = 5
	a[1] = x
	if a[0]+a[1] > 20 {
		return 1
	}
	return 0
}

func (p *PrimitiveArrays) LongArray(a []int64, x int64) int {
	if len(a) != 2 {
		return 255
	}
	a[0] = 5
	a[1] = x
	if a[0]+a[1] > 20 {
		return 1
	}
	return 0
}

func (p *PrimitiveArrays) FloatArray(a []float32, x float32) int {
	if len(a) != 2 {
		return 255
	}
	a[0] = 5
	a[1] = x
	if a[0]+a[1] > 20 {
		return 1
	}
	return 0
}

func (p *PrimitiveArrays) DoubleArray(a []float64, x float64) int {
	if len(a) != 2 {
		return 255
	}
	a[0] = 5
	a[1] = x
	if a[0]+a[1] > 20 {
		return 1
	}
	return 0
}

func (p *PrimitiveArrays) BooleanArray(a []bool, x, y bool) int {
	if len(a) != 2 {
		return 255
	}
	a[0] = x
	a[1] = y
	if a[0] != a[1] {
		return 1
	}
	return 0
}

func (p *PrimitiveArrays) ByteSizeAndIndex(a []byte, x byte) byte {
	if a == nil || len(a) <= int(x) || x < 1 {
		return 255
	}
	b := make([]byte, x)
	b[0] = 5
	a[x] = x
	if b[0]+a[x] > 7 {
		return 1
	}
	return 0
}

func (p *PrimitiveArrays) ShortSizeAndIndex(a []int16, x int16) byte {
	if a == nil || len(a) <= int(x) || x < 1 {
		return 255
	}
	b := make([]int16, x)
	b[0] = 5
	a[x] = x
	if b[0]+a[x] > 7 {
		return 1
	}
	return 0
}

func (p *PrimitiveArrays) CharSizeAndIndex(a []rune, x rune) byte {
	if a == nil || len(a) <= int(x) || x < 1 {
		return 255
	}
	b := make([]rune, x)
	b[0] = 5
	a[x] = x
	if b[0]+a[x] > 7 {
		return 1
	}
	return 0
}
