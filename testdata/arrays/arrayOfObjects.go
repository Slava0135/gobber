package main

import "errors"

type ObjectWithPrimitivesClass struct {
	ValueByDefault int
	x, y           int
	Weight         float64
}

func NewObjectWithPrimitivesClass() *ObjectWithPrimitivesClass {
	return &ObjectWithPrimitivesClass{ValueByDefault: 5}
}

func NewObjectWithPrimitivesClassWithParams(x, y int, weight float64) *ObjectWithPrimitivesClass {
	return &ObjectWithPrimitivesClass{
		ValueByDefault: 5,
		x:              x,
		y:              y,
		Weight:         weight,
	}
}

type ArrayOfObjects struct{}

func (a *ArrayOfObjects) DefaultValues() []*ObjectWithPrimitivesClass {
	array := make([]*ObjectWithPrimitivesClass, 1)
	if array[0] != nil {
		return array
	}
	return array
}

func (a *ArrayOfObjects) CreateArray(x, y, length int) ([]*ObjectWithPrimitivesClass, error) {
	if length < 3 {
		return nil, errors.New("IllegalArgumentException: length must be at least 3")
	}

	array := make([]*ObjectWithPrimitivesClass, length)

	for i := 0; i < len(array); i++ {
		array[i] = NewObjectWithPrimitivesClass()
		array[i].x = x + i
		array[i].y = y + i
	}

	return array, nil
}

func (a *ArrayOfObjects) CopyArray(array []*ObjectWithPrimitivesClass) ([]*ObjectWithPrimitivesClass, error) {
	if len(array) < 3 {
		return nil, errors.New("IllegalArgumentException: array length must be at least 3")
	}

	copy := make([]*ObjectWithPrimitivesClass, len(array))
	for i := 0; i < len(array); i++ {
		copy[i] = array[i]
	}

	for i := 0; i < len(array); i++ {
		array[i].x = -1
		array[i].y = 1
	}

	return copy, nil
}

func (a *ArrayOfObjects) ObjectArray(array []*ObjectWithPrimitivesClass, objectExample *ObjectWithPrimitivesClass) int {
	if len(array) != 2 {
		return -1
	}
	array[0] = NewObjectWithPrimitivesClass()
	array[0].x = 5
	array[1] = objectExample
	if array[0].x+array[1].x > 20 {
		return 1
	}
	return 0
}

type ObjectWithPrimitivesClassSucc struct {
	ObjectWithPrimitivesClass
	AnotherX int
}
