package graph

import "github.com/aclements/go-z3/z3"

type SymValue interface {
	Sort() z3.Sort
}

type Complex struct {
	real z3.Float
	imag z3.Float
	sort z3.Sort
}

type String struct {
	sort z3.Sort
}

func (c *Complex) Sort() z3.Sort {
	return c.sort
}

func (s *String) Sort() z3.Sort {
	return s.sort
}
