package main

import (
	"errors"
	"math"
)

type InvokeClass struct {
	Value int
}

func (i *InvokeClass) DivBy(den int) int {
	return i.Value / den
}

func (i *InvokeClass) UpdateValue(newValue int) {
	i.Value = newValue
}

type InvokeExample struct{}

func (i *InvokeExample) mult(a, b int) int {
	return a * b
}

func (i *InvokeExample) half(a int) int {
	return a / 2
}

func (i *InvokeExample) SimpleFormula(fst, snd int) (int, error) {
	if fst < 100 {
		return 0, errors.New("IllegalArgumentException: fst < 100")
	} else if snd < 100 {
		return 0, errors.New("IllegalArgumentException: snd < 100")
	}

	x := fst + 5
	y := i.half(snd)

	return i.mult(x, y), nil
}

func (i *InvokeExample) initialize(value int) *InvokeClass {
	objectValue := &InvokeClass{
		Value: value,
	}
	return objectValue
}

func (i *InvokeExample) CreateObjectFromValue(value int) *InvokeClass {
	if value == 0 {
		value = 1
	}
	return i.initialize(value)
}

func (i *InvokeExample) changeValue(objectValue *InvokeClass, value int) {
	objectValue.Value = value
}

func (i *InvokeExample) ChangeObjectValueByMethod(objectValue *InvokeClass) *InvokeClass {
	objectValue.Value = 1
	i.changeValue(objectValue, 4)
	return objectValue
}

func (i *InvokeExample) getFive() int {
	return 5
}

func (i *InvokeExample) getTwo() int {
	return 2
}

func (i *InvokeExample) ParticularValue(invokeObject *InvokeClass) (*InvokeClass, error) {
	if invokeObject.Value < 0 {
		return nil, errors.New("IllegalArgumentException: value < 0")
	}
	x := i.getFive() * i.getTwo()
	y := i.getFive() / i.getTwo()

	invokeObject.Value = x + y
	return invokeObject, nil
}

func (i *InvokeExample) getNull() *InvokeClass {
	return nil
}

func (i *InvokeExample) GetNullOrValue(invokeObject *InvokeClass) *InvokeClass {
	if invokeObject.Value < 100 {
		return i.getNull()
	}
	invokeObject.Value = i.getFive()
	return invokeObject
}

func (i *InvokeExample) abs(value int) int {
	if value < 0 {
		if value == math.MinInt {
			return 0
		}
		return i.mult(-1, value)
	}
	return value
}

func (i *InvokeExample) ConstraintsFromOutside(value int) (int, error) {
	if i.abs(value) < 0 {
		return 0, errors.New("IllegalArgumentException: abs(value) < 0")
	}
	return i.abs(value), nil
}

func (i *InvokeExample) helper(value int) int {
	if value < 0 {
		return -1
	}
	return 1
}

func (i *InvokeExample) ConstraintsFromInside(value int) int {
	if value < 0 {
		if value == math.MinInt {
			value = 0
		} else {
			value = -value
		}
	}
	return i.helper(value)
}

func (i *InvokeExample) AlwaysNPE(invokeObject *InvokeClass) *InvokeClass {
	if invokeObject.Value == 0 {
		invokeObject = i.getNull()
		invokeObject.Value = 0
		return invokeObject
	} else if invokeObject.Value > 0 {
		invokeObject = i.getNull()
		invokeObject.Value = 1
		return invokeObject
	} else {
		invokeObject = i.getNull()
		invokeObject.Value = -1
		return invokeObject
	}
}

func (i *InvokeExample) NestedMethodWithException(invokeObject *InvokeClass) (*InvokeClass, error) {
	if invokeObject.Value < 0 {
		return nil, errors.New("IllegalArgumentException: value < 0")
	}
	return invokeObject, nil
}

func (i *InvokeExample) ExceptionInNestedMethod(invokeObject *InvokeClass, value int) (*InvokeClass, error) {
	invokeObject.Value = value
	return i.NestedMethodWithException(invokeObject)
}

func (i *InvokeExample) FewNestedException(invokeObject *InvokeClass, value int) (*InvokeClass, error) {
	invokeObject.Value = value
	return i.firstLevelWithException(invokeObject)
}

func (i *InvokeExample) firstLevelWithException(invokeObject *InvokeClass) (*InvokeClass, error) {
	if invokeObject.Value < 10 {
		return nil, errors.New("IllegalArgumentException: value < 10")
	}
	return i.secondLevelWithException(invokeObject)
}

func (i *InvokeExample) secondLevelWithException(invokeObject *InvokeClass) (*InvokeClass, error) {
	if invokeObject.Value < 100 {
		return nil, errors.New("IllegalArgumentException: value < 100")
	}
	return i.thirdLevelWithException(invokeObject)
}

func (i *InvokeExample) thirdLevelWithException(invokeObject *InvokeClass) (*InvokeClass, error) {
	if invokeObject.Value < 10000 {
		return nil, errors.New("IllegalArgumentException: value < 10000")
	}
	return invokeObject, nil
}

func (i *InvokeExample) DivBy(invokeObject *InvokeClass, den int) (int, error) {
	if invokeObject.Value < 1000 {
		return 0, errors.New("IllegalArgumentException: value < 1000")
	}
	return invokeObject.DivBy(den), nil
}

func (i *InvokeExample) UpdateValue(invokeObject *InvokeClass, value int) (*InvokeClass, error) {
	if invokeObject.Value > 0 {
		return invokeObject, nil
	}
	if value > 0 {
		invokeObject.UpdateValue(value)
		if invokeObject.Value != value {
			return nil, errors.New("RuntimeException: unreachable branch")
		}
		return invokeObject, nil
	}
	return nil, errors.New("IllegalArgumentException")
}

func (i *InvokeExample) NullAsParameter(den int) (int, error) {
	return i.DivBy(nil, den)
}

func (i *InvokeExample) ChangeArrayWithAssignFromMethod(array []int) []int {
	return i.changeAndReturnArray(array, 5)
}

func (i *InvokeExample) changeAndReturnArray(array []int, diff int) []int {
	updatedArray := make([]int, len(array))
	for idx := range array {
		updatedArray[idx] = array[idx] + diff
	}
	return updatedArray
}

func (i *InvokeExample) ChangeArrayByMethod(array []int) []int {
	i.changeArrayValues(array, 5)
	return array
}

func (i *InvokeExample) changeArrayValues(array []int, diff int) {
	for idx := range array {
		array[idx] += diff
	}
}

func (i *InvokeExample) ArrayCopyExample(array []int) ([]int, error) {
	if len(array) < 3 {
		return nil, errors.New("IllegalArgumentException: array length < 3")
	}

	if array[0] <= array[1] || array[1] <= array[2] {
		return nil, nil
	}

	dst := make([]int, len(array))
	i.arrayCopy(array, 0, dst, 0, len(array))

	return dst, nil
}

func (i *InvokeExample) arrayCopy(src []int, srcPos int, dst []int, dstPos int, length int) {
	for j := 0; j < length; j++ {
		dst[dstPos+j] = src[srcPos+j]
	}
}

func (i *InvokeExample) UpdateValues(fst *InvokeClass, snd *InvokeClass) (int, error) {
	i.changeTwoObjects(fst, snd)
	if fst.Value == 1 && snd.Value == 2 {
		return 1, nil
	}
	return 0, errors.New("RuntimeException: values do not match expected output")
}

func (i *InvokeExample) changeTwoObjects(fst *InvokeClass, snd *InvokeClass) {
	fst.Value = 1
	snd.Value = 2
}
