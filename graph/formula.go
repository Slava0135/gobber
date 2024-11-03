package graph

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Formula interface {
	fmt.Stringer
}

type Var struct {
	Name string
	Type string
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

type UnOp struct {
	Result Var
	Arg    Var
	Op     Op
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

type Function struct {
	Result Var
	Name   string
	Args   []Var
}

func (v Var) String() string {
	return v.Name + ":" + v.Type
}

func (o Op) String() string {
	return o.Name
}

func (bo BinOp) String() string {
	return fmt.Sprintf("%s == (%s %s %s)", bo.Result, bo.Left, bo.Op, bo.Right)
}

func (uo UnOp) String() string {
	return fmt.Sprintf("%s == %s%s", uo.Result, uo.Op, uo.Arg)
}

func (ret Return) String() string {
	var s []string
	for _, r := range ret.Results {
		s = append(s, r.String())
	}
	return fmt.Sprintf("return %s", strings.Join(s, ","))
}

func (and And) String() string {
	var s []string
	for _, subf := range and.SubFormulas {
		s = append(s, subf.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(s, ") && ("))
}

func (i If) String() string {
	return fmt.Sprintf("(%s && %s) || (!%s && %s)", i.Cond, i.Then, i.Cond, i.Else)
}

func (f Function) String() string {
	var s []string
	for _, a := range f.Args {
		s = append(s, a.String())
	}
	return fmt.Sprintf("%s == %s(%s)", f.Result, f.Name, strings.Join(s, ", "))
}

func toYaml(f Formula) string {
	d, err := yaml.Marshal(&f)
	if err != nil {
		panic(err)
	}
	return string(d)
}
