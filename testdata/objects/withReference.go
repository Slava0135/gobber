package main

import "errors"

type SimpleDataClass struct {
	a int
	b int
}

type ObjectWithRefFieldClass struct {
	x, y       int
	Weight     float64
	RefField   *SimpleDataClass
	ArrayField []int
}

type ObjectWithRefFieldExample struct{}

func (o *ObjectWithRefFieldExample) DefaultValue() *ObjectWithRefFieldClass {
	obj := &ObjectWithRefFieldClass{}
	if obj.x != 0 || obj.y != 0 || obj.Weight != 0.0 || obj.RefField != nil || obj.ArrayField != nil {
		return obj
	}
	return obj
}

func (o *ObjectWithRefFieldExample) WriteToRefTypeField(objectExample *ObjectWithRefFieldClass, value int) (*ObjectWithRefFieldClass, error) {
	if value != 42 {
		return nil, errors.New("IllegalArgumentException: value must be 42")
	} else if objectExample.RefField != nil {
		return nil, errors.New("IllegalArgumentException: refField must be nil")
	}

	simpleDataClass := &SimpleDataClass{
		a: value,
		b: value * 2,
	}
	objectExample.RefField = simpleDataClass
	return objectExample, nil
}

func (o *ObjectWithRefFieldExample) DefaultFieldValues() *ObjectWithRefFieldClass {
	return &ObjectWithRefFieldClass{}
}

func (o *ObjectWithRefFieldExample) ReadFromRefTypeField(objectExample *ObjectWithRefFieldClass) int {
	if objectExample.RefField == nil || objectExample.RefField.a <= 0 {
		return -1
	}
	return objectExample.RefField.a
}

func (o *ObjectWithRefFieldExample) WriteToArrayField(objectExample *ObjectWithRefFieldClass, length int) (*ObjectWithRefFieldClass, error) {
	if length < 3 {
		return nil, errors.New("IllegalArgumentException: length must be at least 3")
	}

	array := make([]int, length)
	for i := 0; i < length; i++ {
		array[i] = i + 1
	}

	objectExample.ArrayField = array
	objectExample.ArrayField[length-1] = 100

	return objectExample, nil
}

func (o *ObjectWithRefFieldExample) ReadFromArrayField(objectExample *ObjectWithRefFieldClass, value int) int {
	if len(objectExample.ArrayField) > 2 && objectExample.ArrayField[2] == value {
		return 1
	}
	return 2
}

func (o *ObjectWithRefFieldExample) CompareTwoDifferentObjectsFromArguments(fst, snd *ObjectWithRefFieldClass) int {
	if fst.x > 0 && snd.x < 0 {
		if fst == snd {
			panic("RuntimeException")
		} else {
			return 1
		}
	}

	fst.x = snd.x
	fst.y = snd.y
	fst.Weight = snd.Weight

	if fst == snd {
		return 2
	}

	return 3
}

func (o *ObjectWithRefFieldExample) CompareTwoObjectsWithNullRefField(fst, snd *ObjectWithRefFieldClass) int {
	fst.RefField = nil
	snd.RefField = &SimpleDataClass{}
	if fst == snd {
		return 1
	}
	return 2
}

func (o *ObjectWithRefFieldExample) CompareTwoObjectsWithDifferentRefField(fst, snd *ObjectWithRefFieldClass, value int) int {
	fst.RefField = &SimpleDataClass{a: value}
	snd.RefField = &SimpleDataClass{a: fst.RefField.a + 1}
	fst.RefField.b = snd.RefField.b

	if fst == snd {
		return 1
	}
	return 2
}

func (o *ObjectWithRefFieldExample) CompareTwoObjectsWithTheDifferentRefField(fst, snd *ObjectWithRefFieldClass) bool {
	return fst.RefField == snd.RefField
}

func (o *ObjectWithRefFieldExample) CompareTwoObjectsWithTheSameRefField(fst, snd *ObjectWithRefFieldClass) int {
	simpleDataClass := &SimpleDataClass{}

	fst.RefField = simpleDataClass
	snd.RefField = simpleDataClass
	fst.x = snd.x
	fst.y = snd.y
	fst.Weight = snd.Weight

	if fst == snd {
		return 1
	}
	return 2
}
