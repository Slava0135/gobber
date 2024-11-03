package graph

import (
	"fmt"
	"strings"
)

type Formula interface {
	fmt.Stringer
}

type Var struct {
	Name string
}

type Op struct {
	Name string
}

type BinOp struct {
	Result Var
	Left   Var
	Op     Op
	Right  Var
}

type Return struct {
	Results []Var
}

type And struct {
	SubFormulas []Formula
}

type If struct {
	Cond Var
	Then Formula
	Else Formula
}

func (v Var) String() string {
	return v.Name;
}

func (o Op) String() string {
	return o.Name;
}

func (bo BinOp) String() string {
	return fmt.Sprintf("%s == (%s %s %s)", bo.Result, bo.Left, bo.Op, bo.Right);
}

func (ret Return) String() string {
	var s []string
	for _, r := range ret.Results {
		s = append(s, r.String())
	}
	return fmt.Sprintf("return %s", strings.Join(s, ","));
}

func (and And) String() string {
	var s []string
	for _, subf := range and.SubFormulas {
		s = append(s, subf.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(s, ") && ("));
}

func (i If) String() string {
	return fmt.Sprintf("(%s && %s) || (!%s && %s)", i.Cond, i.Then, i.Cond, i.Else);
}
