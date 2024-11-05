package graph

import (
	"fmt"
	"strings"

	"github.com/aclements/go-z3/z3"
	"gopkg.in/yaml.v3"
)

const (
	intType = "int"

	resultSpecialVar = "$result"
)

type Formula interface {
	fmt.Stringer

	Encode(vars map[string]z3.Value, funcs map[string]z3.FuncDecl) z3.Value
}

type Var struct {
	Name     string
	Type     string
	Constant bool
}

type BinOp struct {
	Result Var
	Left   Var
	Op     string
	Right  Var
}

type UnOp struct {
	Result Var
	Arg    Var
	Op     string
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

type Call struct {
	Result Var
	Name   string
	Args   []Var
}

type Convert struct {
	Result Var
	Arg    Var
}

func (v Var) String() string {
	return v.Name + ":" + v.Type
}

func (v Var) Encode(vars map[string]z3.Value, funcs map[string]z3.FuncDecl) z3.Value {
	if v, ok := vars[v.Name]; ok {
		return v
	}
	panic(fmt.Sprintf("unknown var '%s'", v.Name))
}

func (bo BinOp) String() string {
	return fmt.Sprintf("%s == (%s %s %s)", bo.Result, bo.Left, bo.Op, bo.Right)
}

func (bo BinOp) Encode(vars map[string]z3.Value, funcs map[string]z3.FuncDecl) z3.Value {
	res := bo.Result.Encode(vars, funcs)
	left := bo.Left.Encode(vars, funcs)
	right := bo.Right.Encode(vars, funcs)
	switch bo.Op {
	case "+":
		switch res := res.(type) {
		case z3.Int:
			return res.Eq(left.(z3.Int).Add(right.(z3.Int)))
		default:
			panic(fmt.Sprintf("unknown binary operation '%s' for sort '%s'", bo.Op, res.Sort()))
		}
	default:
		panic(fmt.Sprintf("unknown binary operation '%s'", bo.Op))
	}
}

func (uo UnOp) String() string {
	return fmt.Sprintf("%s == %s%s", uo.Result, uo.Op, uo.Arg)
}

func (uo UnOp) Encode(vars map[string]z3.Value, funcs map[string]z3.FuncDecl) z3.Value {
	_ = uo.Result.Encode(vars, funcs)
	_ = uo.Arg.Encode(vars, funcs)
	switch uo.Op {
	default:
		panic(fmt.Sprintf("unknown unary operation '%s'", uo.Op))
	}
}

func (ret Return) String() string {
	var s []string
	for _, r := range ret.Results {
		s = append(s, r.String())
	}
	return fmt.Sprintf("return %s", strings.Join(s, ","))
}

func (ret Return) Encode(vars map[string]z3.Value, funcs map[string]z3.FuncDecl) z3.Value {
	if len(ret.Results) > 1 {
		panic("multiple return values are not supported")
	}
	if result, ok := vars[resultSpecialVar]; ok {
		return result
	}
	panic("result var not found")
}

func (and And) String() string {
	var s []string
	for _, subf := range and.SubFormulas {
		s = append(s, subf.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(s, ") && ("))
}

func (and And) Encode(vars map[string]z3.Value, funcs map[string]z3.FuncDecl) z3.Value {
	var res = and.SubFormulas[0].Encode(vars, funcs).(z3.Bool)
	for i := 1; i < len(and.SubFormulas); i++ {
		res = res.And(and.SubFormulas[i].Encode(vars, funcs).(z3.Bool))
	}
	return res
}

func (i If) String() string {
	return fmt.Sprintf("(%s && %s) || (!%s && %s)", i.Cond, i.Then, i.Cond, i.Else)
}

func (i If) Encode(vars map[string]z3.Value, funcs map[string]z3.FuncDecl) z3.Value {
	var cond = i.Cond.Encode(vars, funcs).(z3.Bool)
	var thn = i.Then.Encode(vars, funcs).(z3.Bool)
	var els = i.Else.Encode(vars, funcs).(z3.Bool)
	return cond.And(thn).Or(cond.Not().And(els))
}

func (f Call) String() string {
	var s []string
	for _, a := range f.Args {
		s = append(s, a.String())
	}
	return fmt.Sprintf("%s == %s(%s)", f.Result, f.Name, strings.Join(s, ", "))
}

func (f Call) Encode(vars map[string]z3.Value, funcs map[string]z3.FuncDecl) z3.Value {
	if fd, ok := funcs[f.Name]; ok {
		var args []z3.Value
		for _, a := range f.Args {
			args = append(args, a.Encode(vars, funcs))
		}
		var res = f.Result.Encode(vars, funcs)
		switch res := res.(type) {
		case z3.Int:
			res.Eq(fd.Apply(args...).(z3.Int))
		default:
			panic(fmt.Sprintf("unknown sort '%s'", res.Sort()))
		}
	}
	panic(fmt.Sprintf("unknown function '%s'", f.Name))
}

func (c Convert) String() string {
	return fmt.Sprintf("%s as %s", c.Arg, c.Result)
}

func (c Convert) Encode(vars map[string]z3.Value, funcs map[string]z3.FuncDecl) z3.Value {
	if c.Result.Type != c.Arg.Type {
		panic("conversions between types are not supported")
	}
	switch c.Result.Type {
	case intType:
		return c.Result.Encode(vars, funcs).(z3.Int).Eq(c.Arg.Encode(vars, funcs).(z3.Int))
	default:
		panic(fmt.Sprintf("unsupported type '%s'", c.Result.Type))
	}
}

func toYaml(f Formula) string {
	d, err := yaml.Marshal(&f)
	if err != nil {
		panic(err)
	}
	return string(d)
}
