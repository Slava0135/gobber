package main

type Point struct{}

type ColoredPoint struct {
	Point
}

type ArrayOfArrays struct{}

func (a *ArrayOfArrays) DefaultValues() [][][]int {
	array := make([][][]int, 1)
	if array[0] == nil {
		return array
	}
	return array
}

func (a *ArrayOfArrays) SizesWithoutTouchingTheElements() [][]int {
	array := make([][]int, 10)
	for i := range array {
		array[i] = make([]int, 3)
	}
	for i := 0; i < 2; i++ {
		_ = array[i][0]
	}
	return array
}

func (a *ArrayOfArrays) DefaultValuesWithoutLastDimension() [][][][]int {
	array := make([][][][]int, 4)
	for i := range array {
		array[i] = make([][][]int, 4)
		for j := range array[i] {
			array[i][j] = make([][]int, 4)
		}
	}
	if array[0][0][0] != nil {
		return array
	}
	return array
}

func (a *ArrayOfArrays) CreateNewMultiDimensionalArray(i, j int) [][][][]int {
	array := make([][][][]int, 2)
	for idx := range array {
		array[idx] = make([][][]int, i)
		for jdx := range array[idx] {
			array[idx][jdx] = make([][]int, j)
			for kdx := range array[idx][jdx] {
				array[idx][jdx][kdx] = make([]int, 3)
			}
		}
	}
	return array
}

func (a *ArrayOfArrays) DefaultValuesWithoutTwoDimensions(i int) [][][][]int {
	if i < 2 {
		return nil
	}
	array := make([][][][]int, i)
	for idx := range array {
		array[idx] = make([][][]int, i)
	}
	if array[0][0] != nil {
		return array
	}
	return array
}

func (a *ArrayOfArrays) DefaultValuesNewMultiArray() [][][]int {
	array := make([][][]int, 1)
	array[0] = make([][]int, 1)
	array[0][0] = make([]int, 1)
	if array[0] == nil {
		return array
	}
	return array
}

func (a *ArrayOfArrays) IsIdentityMatrix(matrix [][]int) bool {
	if len(matrix) < 3 {
		panic("IllegalArgumentException: matrix length < 3")
	}
	for i := 0; i < len(matrix); i++ {
		if len(matrix[i]) != len(matrix) {
			return false
		}
		for j := 0; j < len(matrix[i]); j++ {
			if i == j && matrix[i][j] != 1 {
				return false
			}
			if i != j && matrix[i][j] != 0 {
				return false
			}
		}
	}
	return true
}

func (a *ArrayOfArrays) CreateNewThreeDimensionalArray(length, constValue int) [][][]int {
	if length != 2 {
		return make([][][]int, 0)
	}
	matrix := make([][][]int, length)
	for i := range matrix {
		matrix[i] = make([][]int, length)
		for j := range matrix[i] {
			matrix[i][j] = make([]int, length)
			for k := range matrix[i][j] {
				matrix[i][j][k] = constValue + 1
			}
		}
	}
	return matrix
}

func (a *ArrayOfArrays) ReallyMultiDimensionalArray(array [][][]int) [][][]int {
	if array[1][2][3] != 12345 {
		array[1][2][3] = 12345
	} else {
		array[1][2][3] -= 12345 * 2
	}
	return array
}

func (a *ArrayOfArrays) MultiDimensionalObjectsArray(array [][]Point) [][]Point {
	array[0] = make([]Point, 2)
	array[1] = make([]Point, 1)
	return array
}

func (a *ArrayOfArrays) FillMultiArrayWithArray(value []int) [][]int {
	if len(value) < 2 {
		return make([][]int, 0)
	}
	for i := range value {
		value[i] += i
	}
	length := 3
	array := make([][]int, length)
	for i := 0; i < length; i++ {
		array[i] = value
	}
	return array
}
