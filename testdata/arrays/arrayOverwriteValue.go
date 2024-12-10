package main

type ArraysOverwriteValue struct{}

func (a *ArraysOverwriteValue) ByteArray(arr []byte) byte {
	if len(arr) == 0 {
		return 1
	}
	if arr[0] != 0 {
		return 2
	}
	arr[0] = 1
	return 3
}

func (a *ArraysOverwriteValue) ShortArray(arr []int16) byte {
	if len(arr) == 0 {
		return 1
	}
	if arr[0] != 0 {
		return 2
	}
	arr[0] = 1
	return 3
}

func (a *ArraysOverwriteValue) CharArray(arr []rune) rune {
	if len(arr) == 0 {
		return 1
	}
	if arr[0] != 0 {
		return 2
	}
	arr[0] = 1
	return 3
}

func (a *ArraysOverwriteValue) IntArray(arr []int) byte {
	if len(arr) == 0 {
		return 1
	}
	if arr[0] != 0 {
		return 2
	}
	arr[0] = 1
	return 3
}

func (a *ArraysOverwriteValue) LongArray(arr []int64) byte {
	if len(arr) == 0 {
		return 1
	}
	if arr[0] != 0 {
		return 2
	}
	arr[0] = 1
	return 3
}

func (a *ArraysOverwriteValue) FloatArray(arr []float32) byte {
	if len(arr) == 0 {
		return 1
	}
	if arr[0] == arr[0] {
		return 2
	}
	arr[0] = 1.0
	return 3
}

func (a *ArraysOverwriteValue) DoubleArray(arr []float64) byte {
	if len(arr) == 0 {
		return 1
	}
	if arr[0] == arr[0] {
		return 2
	}
	arr[0] = 1.0
	return 3
}

func (a *ArraysOverwriteValue) BooleanArray(arr []bool) int {
	if len(arr) == 0 {
		return 1
	}
	if arr[0] {
		return 2
	}
	arr[0] = true
	return 3
}

func (a *ArraysOverwriteValue) ObjectArray(arr []*ObjectWithPrimitivesClass) int {
	if len(arr) == 0 {
		return 1
	}
	if arr[0] == nil {
		return 2
	}
	arr[0] = nil
	return 3
}
