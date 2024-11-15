package symexec

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aclements/go-z3/z3"
	"gopkg.in/yaml.v3"
)

const (
	unsignedIntType = "uint"
	intType         = "int"
	boolType        = "bool"
	floatType       = "float64"
	complexType     = "complex128"
	stringType      = "string"
	intArrayType    = "[]int"
	intPointerType  = "*int"

	intSize     = 64
	floatSize   = 64
	complexSize = 64

	resultSpecialVar = "$result"
)

type EncodingContext struct {
	*z3.Context

	vars  map[string]SymValue
	funcs map[string]z3.FuncDecl

	intValues            z3.Array
	intArrayValuesMemory z3.Array
	intArrayLenMemory    z3.Array

	floatSort      z3.Sort
	complexSort    z3.Sort
	stringSort     z3.Sort
	intArraySort   z3.Sort
	intPointerSort z3.Sort

	addrSort z3.Sort
	null     z3.Sort
}

func (ctx *EncodingContext) ComplexConst(name string) *Complex {
	return &Complex{
		real: ctx.Const(name+".REAL", ctx.floatSort).(z3.Float),
		imag: ctx.Const(name+".IMAG", ctx.floatSort).(z3.Float),
		sort: ctx.complexSort,
	}
}

func (ctx *EncodingContext) StringConst(name string) *String {
	return &String{
		sort: ctx.stringSort,
	}
}

func (ctx *EncodingContext) IntArrayConst(name string) *IntArray {
	return &IntArray{
		addr: ctx.Const(name, ctx.addrSort).(z3.Uninterpreted),
		sort: ctx.intArraySort,
	}
}

func (ctx *EncodingContext) IntPointerConst(name string) *IntPointer {
	return &IntPointer{
		addr: ctx.Const(name, ctx.addrSort).(z3.Uninterpreted),
		sort: ctx.intPointerSort,
	}
}

func (ctx *EncodingContext) FromComplex128(c complex128) *Complex {
	return &Complex{
		real: ctx.FromFloat64(real(c), ctx.floatSort),
		imag: ctx.FromFloat64(imag(c), ctx.floatSort),
		sort: ctx.complexSort,
	}
}

type Formula interface {
	fmt.Stringer

	Encode(ctx *EncodingContext) SymValue
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

type IndexAddr struct {
	Result Var
	Array  Var
	Index  Var
}

func (v Var) String() string {
	return v.Name + ":" + v.Type
}

func (v Var) Encode(ctx *EncodingContext) SymValue {
	if v.Constant {
		switch v.Type {
		case intType, unsignedIntType:
			i, err := strconv.ParseInt(v.Name, 10, intSize)
			if err != nil {
				panic(err)
			}
			return ctx.FromInt(i, ctx.IntSort())
		case boolType:
			b, err := strconv.ParseBool(v.Name)
			if err != nil {
				panic(err)
			}
			return ctx.FromBool(b)
		case floatType:
			f, err := strconv.ParseFloat(v.Name, floatSize)
			if err != nil {
				panic(err)
			}
			return ctx.FromFloat64(f, ctx.floatSort)
		case complexType:
			c, err := strconv.ParseComplex(v.Name, complexSize)
			if err != nil {
				panic(err)
			}
			return ctx.FromComplex128(c)
		default:
			panic(fmt.Sprintf("unknown constant '%s' of type '%s'", v.Name, v.Type))
		}
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

func (bo BinOp) Encode(ctx *EncodingContext) SymValue {
	unknownOp := func(op string, sort z3.Sort) {
		panic(fmt.Sprintf("unknown binary operation '%s' for sort '%s'", op, sort))
	}
	res := bo.Result.Encode(ctx)
	left := bo.Left.Encode(ctx)
	right := bo.Right.Encode(ctx)
	switch bo.Op {
	case "+":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Int).Eq(left.Add(right.(z3.Int)))
		case z3.Float:
			return res.(z3.Float).Eq(left.Add(right.(z3.Float)))
		case *Complex:
			resCx := res.(*Complex)
			rightCx := right.(*Complex)
			return resCx.real.Eq(left.real.Add(rightCx.real)).And(
				resCx.imag.Eq(left.imag.Add(rightCx.imag)))
		default:
			unknownOp(bo.Op, left.Sort())
		}
	case "-":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Int).Eq(left.Sub(right.(z3.Int)))
		case z3.Float:
			return res.(z3.Float).Eq(left.Sub(right.(z3.Float)))
		case *Complex:
			resCx := res.(*Complex)
			rightCx := right.(*Complex)
			return resCx.real.Eq(left.real.Sub(rightCx.real)).And(
				resCx.imag.Eq(left.imag.Sub(rightCx.imag)))
		default:
			unknownOp(bo.Op, left.Sort())
		}
	case "*":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Int).Eq(left.Mul(right.(z3.Int)))
		case z3.Float:
			return res.(z3.Float).Eq(left.Mul(right.(z3.Float)))
		case *Complex:
			resCx := res.(*Complex)
			rightCx := right.(*Complex)
			// (a+bi)(c+di) = (ac - bd) + (ad + bc)i
			a := left.real
			b := left.imag
			c := rightCx.real
			d := rightCx.imag
			return resCx.real.Eq(a.Mul(c).Sub(b.Mul(d))).And(
				resCx.imag.Eq(a.Mul(d).Add(b.Mul(c))))
		default:
			unknownOp(bo.Op, left.Sort())
		}
	case "/":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Int).Eq(left.Div(right.(z3.Int)))
		case z3.Float:
			return res.(z3.Float).Eq(left.Div(right.(z3.Float)))
		case *Complex:
			resCx := res.(*Complex)
			rightCx := right.(*Complex)
			// (a+bi)/(c+di) = ((ac + bd) + (bc - ad)i)/(c^2 + d^2)
			a := left.real
			b := left.imag
			c := rightCx.real
			d := rightCx.imag
			denom := c.Mul(c).Add(d.Mul(d))
			return resCx.real.Eq(a.Mul(c).Add(b.Mul(d)).Div(denom)).And(
				resCx.imag.Eq(b.Mul(c).Sub(a.Mul(d)).Div(denom)))
		default:
			unknownOp(bo.Op, left.Sort())
		}
	case "%":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Int).Eq(left.Mod(right.(z3.Int)))
		default:
			unknownOp(bo.Op, left.Sort())
		}
	case ">":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Bool).Eq(left.GT(right.(z3.Int)))
		case z3.Float:
			return res.(z3.Bool).Eq(left.GT(right.(z3.Float)))
		default:
			unknownOp(bo.Op, left.Sort())
		}
	case ">=":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Bool).Eq(left.GE(right.(z3.Int)))
		case z3.Float:
			return res.(z3.Bool).Eq(left.GE(right.(z3.Float)))
		default:
			unknownOp(bo.Op, left.Sort())
		}
	case "<":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Bool).Eq(left.LT(right.(z3.Int)))
		case z3.Float:
			return res.(z3.Bool).Eq(left.LT(right.(z3.Float)))
		default:
			unknownOp(bo.Op, left.Sort())
		}
	case "<=":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Bool).Eq(left.LE(right.(z3.Int)))
		case z3.Float:
			return res.(z3.Bool).Eq(left.LE(right.(z3.Float)))
		default:
			unknownOp(bo.Op, left.Sort())
		}
	case "==":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Bool).Eq(left.Eq(right.(z3.Int)))
		case z3.Float:
			return res.(z3.Bool).Eq(left.IEEEEq(right.(z3.Float)))
		default:
			unknownOp(bo.Op, left.Sort())
		}
	case "<<":
		switch left := left.(type) {
		case z3.Int:
			leftBV := left.ToBV(intSize)
			rightBV := right.(z3.Int).ToBV(intSize)
			return res.(z3.Int).Eq(leftBV.Lsh(rightBV).SToInt())
		default:
			unknownOp(bo.Op, left.Sort())
		}
	case ">>":
		switch left := left.(type) {
		case z3.Int:
			leftBV := left.ToBV(intSize)
			rightBV := right.(z3.Int).ToBV(intSize)
			return res.(z3.Int).Eq(leftBV.SRsh(rightBV).SToInt())
		default:
			unknownOp(bo.Op, left.Sort())
		}
	case "^":
		switch left := left.(type) {
		case z3.Int:
			leftBV := left.ToBV(intSize)
			rightBV := right.(z3.Int).ToBV(intSize)
			return res.(z3.Int).Eq(leftBV.Xor(rightBV).SToInt())
		default:
			unknownOp(bo.Op, left.Sort())
		}
	case "&":
		switch left := left.(type) {
		case z3.Int:
			leftBV := left.ToBV(intSize)
			rightBV := right.(z3.Int).ToBV(intSize)
			return res.(z3.Int).Eq(leftBV.And(rightBV).SToInt())
		default:
			unknownOp(bo.Op, left.Sort())
		}
	case "|":
		switch left := left.(type) {
		case z3.Int:
			leftBV := left.ToBV(intSize)
			rightBV := right.(z3.Int).ToBV(intSize)
			return res.(z3.Int).Eq(leftBV.Or(rightBV).SToInt())
		default:
			unknownOp(bo.Op, left.Sort())
		}
	default:
		panic(fmt.Sprintf("unknown binary operation '%s'", bo.Op))
	}
	panic("unreachable")
}

func (bo BinOp) ScanVars(vars map[string]Var) {
	bo.Result.ScanVars(vars)
	bo.Left.ScanVars(vars)
	bo.Right.ScanVars(vars)
}

func (uo UnOp) String() string {
	return fmt.Sprintf("%s == %s%s", uo.Result, uo.Op, uo.Arg)
}

func (uo UnOp) Encode(ctx *EncodingContext) SymValue {
	_ = uo.Result.Encode(ctx)
	_ = uo.Arg.Encode(ctx)
	switch uo.Op {
	case "*":
		res := uo.Result.Encode(ctx).(z3.Int)
		arg := uo.Arg.Encode(ctx).(*IntPointer)
		return res.Eq(ctx.intValues.Select(arg.addr).(z3.Int))
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

func (ret Return) Encode(ctx *EncodingContext) SymValue {
	if len(ret.Results) > 1 {
		panic("multiple return values are not supported")
	}
	if result, ok := ctx.vars[resultSpecialVar]; ok {
		switch result := result.(type) {
		case z3.Int:
			return result.Eq(ret.Results[0].Encode(ctx).(z3.Int))
		case z3.Bool:
			return result.Eq(ret.Results[0].Encode(ctx).(z3.Bool))
		case z3.Float:
			return result.Eq(ret.Results[0].Encode(ctx).(z3.Float))
		case *Complex:
			arg := ret.Results[0].Encode(ctx).(*Complex)
			return result.real.Eq(arg.real).And(result.imag.Eq(arg.imag))
		case *String:
			return result
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

func (and And) Encode(ctx *EncodingContext) SymValue {
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

func (i If) Encode(ctx *EncodingContext) SymValue {
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

func (f Call) Encode(ctx *EncodingContext) SymValue {
	// built-in
	switch f.Name {
	case "real":
		return f.Result.Encode(ctx).(z3.Float).Eq(f.Args[0].Encode(ctx).(*Complex).real)
	case "imag":
		return f.Result.Encode(ctx).(z3.Float).Eq(f.Args[0].Encode(ctx).(*Complex).imag)
	case "len":
		intArray := f.Args[0].Encode(ctx).(*IntArray)
		return f.Result.Encode(ctx).(z3.Int).Eq(ctx.intArrayLenMemory.Select(intArray.addr).(z3.Int))
	}
	if function, ok := ctx.funcs[f.Name]; ok {
		var args []z3.Value
		for _, a := range f.Args {
			args = append(args, a.Encode(ctx).(z3.Value))
		}
		var res = f.Result.Encode(ctx)
		switch res := res.(type) {
		case z3.Int:
			res.Eq(function.Apply(args...).(z3.Int))
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

func (c Convert) Encode(ctx *EncodingContext) SymValue {
	unsupportedConv := func() {
		panic(fmt.Sprintf("unsupported conversion from '%s' to '%s'", c.Arg.Type, c.Result.Type))
	}
	switch c.Result.Type {
	case intType, unsignedIntType:
		switch c.Arg.Type {
		case intType, unsignedIntType:
			return c.Result.Encode(ctx).(z3.Int).Eq(c.Arg.Encode(ctx).(z3.Int))
		default:
			unsupportedConv()
		}
	case boolType:
		switch c.Arg.Type {
		case boolType:
			return c.Result.Encode(ctx).(z3.Bool).Eq(c.Arg.Encode(ctx).(z3.Bool))
		default:
			unsupportedConv()
		}
	case floatType:
		switch c.Arg.Type {
		case floatType:
			return c.Result.Encode(ctx).(z3.Float).Eq(c.Arg.Encode(ctx).(z3.Float))
		case intType, unsignedIntType:
			return c.Result.Encode(ctx).(z3.Float).Eq(c.Arg.Encode(ctx).(z3.Int).ToBV(intSize).IEEEToFloat(ctx.floatSort))
		default:
			unsupportedConv()
		}
	case complexType:
		switch c.Arg.Type {
		case complexType:
			res := c.Result.Encode(ctx).(*Complex)
			arg := c.Arg.Encode(ctx).(*Complex)
			return res.real.Eq(arg.real).And(res.imag.Eq(arg.imag))
		}
	default:
		unsupportedConv()
	}
	panic("unreachable")
}

func (c Convert) ScanVars(vars map[string]Var) {
	c.Arg.ScanVars(vars)
	c.Result.ScanVars(vars)
}

func (ia IndexAddr) String() string {
	return fmt.Sprintf("%s = &%s[%s]", ia.Result, ia.Array, ia.Index)
}

func (ia IndexAddr) Encode(ctx *EncodingContext) SymValue {
	res := ia.Result.Encode(ctx).(*IntPointer).addr
	array := ia.Array.Encode(ctx).(*IntArray)
	index := ia.Index.Encode(ctx).(z3.Int)
	values := ctx.intArrayValuesMemory.Select(array.addr).(z3.Array)
	value := values.Select(index).(z3.Uninterpreted)
	return res.Eq(value)
}

func (ia IndexAddr) ScanVars(vars map[string]Var) {
	ia.Result.ScanVars(vars)
	ia.Array.ScanVars(vars)
	ia.Index.ScanVars(vars)
}

func toYaml(f Formula) string {
	d, err := yaml.Marshal(&f)
	if err != nil {
		panic(err)
	}
	return string(d)
}
