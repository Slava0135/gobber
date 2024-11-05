package graph

import (
	"fmt"
	"strings"

	"github.com/aclements/go-z3/z3"
	"gopkg.in/yaml.v3"
)

const (
	intType  = "int"
	boolType = "bool"

	resultSpecialVar = "$result"
)

type EncodingContext struct {
	*z3.Context
	vars  map[string]z3.Value
	funcs map[string]z3.FuncDecl
}

type Formula interface {
	fmt.Stringer

	Encode(ctx *EncodingContext) z3.Value
	ScanVars(vars map[string]Var)
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

func (v Var) Encode(ctx *EncodingContext) z3.Value {
	if v.Constant {
		panic("TODO: constants")
	}
	if v, ok := ctx.vars[v.Name]; ok {
		return v
	}
	panic(fmt.Sprintf("unknown var '%s'", v.Name))
}

func (v Var) ScanVars(vars map[string]Var) {
	if v.Constant {
		return
	}
	if oldV, ok := vars[v.Name]; ok {
		if oldV.Type != v.Type {
			panic(fmt.Sprintf("variable '%s' can't different types ('%s' and '%s')", v.Name, oldV.Type, v.Type))
		}
		return
	}
	vars[v.Name] = v
}

func (bo BinOp) String() string {
	return fmt.Sprintf("%s == (%s %s %s)", bo.Result, bo.Left, bo.Op, bo.Right)
}

func (bo BinOp) Encode(ctx *EncodingContext) z3.Value {
	res := bo.Result.Encode(ctx)
	left := bo.Left.Encode(ctx)
	right := bo.Right.Encode(ctx)
	switch bo.Op {
	case "+":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Int).Eq(left.Add(right.(z3.Int)))
		default:
			panic(fmt.Sprintf("unknown binary operation '%s' for sort '%s'", bo.Op, left.Sort()))
		}
	case "-":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Int).Eq(left.Sub(right.(z3.Int)))
		default:
			panic(fmt.Sprintf("unknown binary operation '%s' for sort '%s'", bo.Op, left.Sort()))
		}
	case "*":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Int).Eq(left.Mul(right.(z3.Int)))
		default:
			panic(fmt.Sprintf("unknown binary operation '%s' for sort '%s'", bo.Op, left.Sort()))
		}
	case ">":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Bool).Eq(left.GT(right.(z3.Int)))
		default:
			panic(fmt.Sprintf("unknown binary operation '%s' for sort '%s'", bo.Op, left.Sort()))
		}
	case "<":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Bool).Eq(left.LT(right.(z3.Int)))
		default:
			panic(fmt.Sprintf("unknown binary operation '%s' for sort '%s'", bo.Op, left.Sort()))
		}
	default:
		panic(fmt.Sprintf("unknown binary operation '%s'", bo.Op))
	}
}

func (bo BinOp) ScanVars(vars map[string]Var) {
	bo.Result.ScanVars(vars)
	bo.Left.ScanVars(vars)
	bo.Right.ScanVars(vars)
}

func (uo UnOp) String() string {
	return fmt.Sprintf("%s == %s%s", uo.Result, uo.Op, uo.Arg)
}

func (uo UnOp) Encode(ctx *EncodingContext) z3.Value {
	_ = uo.Result.Encode(ctx)
	_ = uo.Arg.Encode(ctx)
	switch uo.Op {
	default:
		panic(fmt.Sprintf("unknown unary operation '%s'", uo.Op))
	}
}

func (uo UnOp) ScanVars(vars map[string]Var) {
	uo.Arg.ScanVars(vars)
	uo.Result.ScanVars(vars)
}

func (ret Return) String() string {
	var s []string
	for _, r := range ret.Results {
		s = append(s, r.String())
	}
	return fmt.Sprintf("return %s", strings.Join(s, ","))
}

func (ret Return) Encode(ctx *EncodingContext) z3.Value {
	if len(ret.Results) > 1 {
		panic("multiple return values are not supported")
	}
	if result, ok := ctx.vars[resultSpecialVar]; ok {
		switch result := result.(type) {
		case z3.Int:
			return result.Eq(ctx.vars[ret.Results[0].Name].(z3.Int))
		default:
			panic(fmt.Sprintf("unknown return sort '%s'", result.Sort()))
		}
	}
	panic("result var not found")
}

func (ret Return) ScanVars(vars map[string]Var) {
	if len(ret.Results) > 1 {
		panic("multiple return values are not supported")
	}
	res := ret.Results[0]
	if v, ok := vars[resultSpecialVar]; ok {
		if res.Type != v.Type {
			panic(fmt.Sprintf("return values can't have different types ('%s' and '%s')", res.Type, v.Type))
		}
	} else {
		vars[resultSpecialVar] = Var{
			Name:     resultSpecialVar,
			Type:     res.Type,
			Constant: res.Constant,
		}
	}
	for _, v := range ret.Results {
		v.ScanVars(vars)
	}
}

func (and And) String() string {
	var s []string
	for _, subf := range and.SubFormulas {
		s = append(s, subf.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(s, ") && ("))
}

func (and And) Encode(ctx *EncodingContext) z3.Value {
	var res = and.SubFormulas[0].Encode(ctx).(z3.Bool)
	for i := 1; i < len(and.SubFormulas); i++ {
		res = res.And(and.SubFormulas[i].Encode(ctx).(z3.Bool))
	}
	return res
}

func (and And) ScanVars(vars map[string]Var) {
	for _, f := range and.SubFormulas {
		f.ScanVars(vars)
	}
}

func (i If) String() string {
	return fmt.Sprintf("(%s && %s) || (!%s && %s)", i.Cond, i.Then, i.Cond, i.Else)
}

func (i If) Encode(ctx *EncodingContext) z3.Value {
	var cond = i.Cond.Encode(ctx).(z3.Bool)
	var thn = i.Then.Encode(ctx).(z3.Bool)
	var els = i.Else.Encode(ctx).(z3.Bool)
	return cond.And(thn).Or(cond.Not().And(els))
}

func (i If) ScanVars(vars map[string]Var) {
	i.Cond.ScanVars(vars)
	i.Then.ScanVars(vars)
	i.Else.ScanVars(vars)
}

func (f Call) String() string {
	var s []string
	for _, a := range f.Args {
		s = append(s, a.String())
	}
	return fmt.Sprintf("%s == %s(%s)", f.Result, f.Name, strings.Join(s, ", "))
}

func (f Call) Encode(ctx *EncodingContext) z3.Value {
	if fd, ok := ctx.funcs[f.Name]; ok {
		var args []z3.Value
		for _, a := range f.Args {
			args = append(args, a.Encode(ctx))
		}
		var res = f.Result.Encode(ctx)
		switch res := res.(type) {
		case z3.Int:
			res.Eq(fd.Apply(args...).(z3.Int))
		default:
			panic(fmt.Sprintf("unknown sort '%s'", res.Sort()))
		}
	}
	panic(fmt.Sprintf("unknown function '%s'", f.Name))
}

func (f Call) ScanVars(vars map[string]Var) {
	f.Result.ScanVars(vars)
	for _, a := range f.Args {
		a.ScanVars(vars)
	}
}

func (c Convert) String() string {
	return fmt.Sprintf("%s as %s", c.Arg, c.Result)
}

func (c Convert) Encode(ctx *EncodingContext) z3.Value {
	if c.Result.Type != c.Arg.Type {
		panic("conversions between types are not supported")
	}
	switch c.Result.Type {
	case intType:
		return c.Result.Encode(ctx).(z3.Int).Eq(c.Arg.Encode(ctx).(z3.Int))
	default:
		panic(fmt.Sprintf("unsupported type '%s'", c.Result.Type))
	}
}

func (c Convert) ScanVars(vars map[string]Var) {
	c.Arg.ScanVars(vars)
	c.Result.ScanVars(vars)
}

func toYaml(f Formula) string {
	d, err := yaml.Marshal(&f)
	if err != nil {
		panic(err)
	}
	return string(d)
}
