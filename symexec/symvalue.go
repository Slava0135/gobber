package symexec

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

type Pointer struct {
	addr z3.Uninterpreted
	t    string
	elem string
	sort z3.Sort
}

type SymArray struct {
	addr z3.Uninterpreted
	t    string
	sort z3.Sort
}

type SymStruct struct {
	addr z3.Uninterpreted
	t    string
	sort z3.Sort
}

func (c *Complex) Sort() z3.Sort {
	return c.sort
}

func (s *String) Sort() z3.Sort {
	return s.sort
}

func (p *Pointer) Sort() z3.Sort {
	return p.sort
}

func (sa *SymArray) Sort() z3.Sort {
	return sa.sort
}

func (ss *SymStruct) Sort() z3.Sort {
	return ss.sort
}
