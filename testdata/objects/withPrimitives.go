package main

import (
	"errors"
	"fmt"
)

type ObjectWithPrimitivesClass struct {
	ValueByDefault int
	x, y           int
	ShortValue     int16
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

func (o *ObjectWithPrimitivesClass) Equals(other *ObjectWithPrimitivesClass) bool {
	if other == nil {
		return false
	}
	return o.ValueByDefault == other.ValueByDefault && o.x == other.x && o.y == other.y && o.Weight == other.Weight
}

func (o *ObjectWithPrimitivesClass) String() string {
	return fmt.Sprintf(
		"ObjectWithPrimitivesClass(ValueByDefault = %d, x = %d, y = %d, ShortValue = %d, Weight = %f)",
		o.ValueByDefault,
		o.x,
		o.y,
		o.ShortValue,
		o.Weight,
	)
}

type ObjectWithPrimitivesExample struct{}

func (o *ObjectWithPrimitivesExample) Max(fst, snd *ObjectWithPrimitivesClass) *ObjectWithPrimitivesClass {
	if fst.x > snd.x && fst.y > snd.y {
		return fst
	} else if fst.x < snd.x && fst.y < snd.y {
		return snd
	}
	return fst
}

func (o *ObjectWithPrimitivesExample) IgnoredInputParameters(fst, snd *ObjectWithPrimitivesClass) *ObjectWithPrimitivesClass {
	return NewObjectWithPrimitivesClass()
}

func (o *ObjectWithPrimitivesExample) Example(value *ObjectWithPrimitivesClass) *ObjectWithPrimitivesClass {
	if value.x == 1 {
		return value
	}
	value.x = 1
	return value
}

func (o *ObjectWithPrimitivesExample) DefaultValueForSuperclassFields() *ObjectWithPrimitivesClassSucc {
	obj := &ObjectWithPrimitivesClassSucc{}
	if obj.x != 0 {
		return obj
	}
	return obj
}

func (o *ObjectWithPrimitivesExample) CreateObject(a, b int, objectExample *ObjectWithPrimitivesClass) (*ObjectWithPrimitivesClass, error) {
	object := NewObjectWithPrimitivesClass()
	object.x = a + 5
	object.y = b + 6
	object.Weight = objectExample.Weight
	if object.Weight < 0 {
		return nil, errors.New("IllegalArgumentException: weight < 0")
	}
	return object, nil
}

func (o *ObjectWithPrimitivesExample) Memory(objectExample *ObjectWithPrimitivesClass, value int) *ObjectWithPrimitivesClass {
	if value > 0 {
		objectExample.x = 1
		objectExample.y = 2
		objectExample.Weight = 1.2
	} else {
		objectExample.x = -1
		objectExample.y = -2
		objectExample.Weight = -1.2
	}
	return objectExample
}

func (o *ObjectWithPrimitivesExample) NullExample(object *ObjectWithPrimitivesClass) *ObjectWithPrimitivesClass {
	if object.x == 0 && object.y == 0 {
		return nil
	}
	return object
}

func (o *ObjectWithPrimitivesExample) CompareTwoNullObjects(value int) int {
	fst := NewObjectWithPrimitivesClass()
	snd := NewObjectWithPrimitivesClass()

	fst.x = value + 1
	snd.x = value + 2

	if fst.x == value+1 {
		fst = nil
	}
	if snd.x == value+2 {
		snd = nil
	}

	if fst == snd {
		return 1
	}
	panic("RuntimeException")
}

func (o *ObjectWithPrimitivesExample) CompareTwoOuterObjects(fst, snd *ObjectWithPrimitivesClass) bool {
	if fst == nil || snd == nil {
		panic("NullPointerException")
	}
	return fst == snd
}

func (o *ObjectWithPrimitivesExample) CompareTwoDifferentObjects(a int) int {
	fst := NewObjectWithPrimitivesClass()
	snd := NewObjectWithPrimitivesClass()

	fst.x = a
	snd.x = a

	if fst == snd {
		panic("RuntimeException")
	}
	return 1
}

func (o *ObjectWithPrimitivesExample) CompareTwoRefEqualObjects(a int) int {
	fst := NewObjectWithPrimitivesClass()
	snd := fst

	fst.x = a

	if fst == snd {
		return 1
	}
	panic("RuntimeException")
}

func (o *ObjectWithPrimitivesExample) CompareObjectWithArgument(fst *ObjectWithPrimitivesClass) int {
	snd := NewObjectWithPrimitivesClass()

	if snd == fst {
		panic("RuntimeException")
	}
	return 1
}

func (o *ObjectWithPrimitivesExample) CompareTwoIdenticalObjectsFromArguments(fst, snd *ObjectWithPrimitivesClass) int {
	fst.x = snd.x
	fst.y = snd.y
	fst.Weight = snd.Weight

	if fst == snd {
		return 1
	}
	return 2
}

func (o *ObjectWithPrimitivesExample) GetOrDefault(objectExample, defaultValue *ObjectWithPrimitivesClass) (*ObjectWithPrimitivesClass, error) {
	if defaultValue.x == 0 && defaultValue.y == 0 {
		return nil, errors.New("IllegalArgumentException")
	}
	if objectExample == nil {
		return defaultValue, nil
	}
	return objectExample, nil
}

func (o *ObjectWithPrimitivesExample) InheritorsFields(fst *ObjectWithPrimitivesClassSucc, snd *ObjectWithPrimitivesClass) int {
	fst.x = 1
	fst.AnotherX = 2
	fst.y = 3
	fst.Weight = 4.5

	snd.x = 1
	snd.y = 2
	snd.Weight = 3.4

	return 1
}

func (o *ObjectWithPrimitivesExample) CreateWithConstructor(x, y int) *ObjectWithPrimitivesClass {
	return NewObjectWithPrimitivesClassWithParams(x+1, y+2, 3.3)
}

func (o *ObjectWithPrimitivesExample) CreateWithSuperConstructor(x, y, anotherX int) *ObjectWithPrimitivesClassSucc {
	return &ObjectWithPrimitivesClassSucc{
		ObjectWithPrimitivesClass: ObjectWithPrimitivesClass{
			ValueByDefault: 5,
			x:              x + 1,
			y:              y + 2,
			Weight:         3.3,
		},
		AnotherX: anotherX + 4,
	}
}

func (o *ObjectWithPrimitivesExample) FieldWithDefaultValue(x, y int) *ObjectWithPrimitivesClass {
	return NewObjectWithPrimitivesClassWithParams(x, y, 3.3)
}

func (o *ObjectWithPrimitivesExample) ValueByDefault() int {
	objectExample := NewObjectWithPrimitivesClass()
	return objectExample.ValueByDefault
}

type ObjectWithPrimitivesClassSucc struct {
	ObjectWithPrimitivesClass
	AnotherX int
}
