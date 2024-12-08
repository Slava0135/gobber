package main

import "errors"

type RecursiveTypeClass struct {
	Next  *RecursiveTypeClass
	Value int
}

type RecursiveType struct{}

func (r *RecursiveType) NextValue(node *RecursiveTypeClass, value int) (*RecursiveTypeClass, error) {
	if value == 0 {
		return nil, errors.New("IllegalArgumentException: value cannot be 0")
	}
	if node.Next != nil && node.Next.Value == value {
		return node.Next, nil
	}
	return nil, nil
}

func (r *RecursiveType) WriteObjectField(node *RecursiveTypeClass) *RecursiveTypeClass {
	if node.Next == nil {
		node.Next = &RecursiveTypeClass{}
	}
	node.Next.Value = node.Next.Value + 1
	return node
}
