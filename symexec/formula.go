package symexec

import (
	"fmt"
	"go/types"
	"math/big"
	"strconv"
	"strings"

	"github.com/aclements/go-z3/z3"
	"gopkg.in/yaml.v3"
)

type Formula interface {
	fmt.Stringer

	Encode(ctx *EncodingContext) SymValue
	ScanVars(vars map[string]Var)
}

type Var struct {
	Name     string
	Type     types.Type
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

type FieldAddr struct {
	Result Var
	Object Var
	Field  int
}

func (v Var) String() string {
	return v.Name + ":" + v.Type.String()
}

func (v Var) Encode(ctx *EncodingContext) SymValue {
	if v.Constant {
		switch t := v.Type.(type) {
		case *types.Basic:
			switch t.Kind() {
			case types.Uint:
				i, err := strconv.ParseUint(v.Name, 10, intSize)
				if err != nil {
					panic(err)
				}
				return ctx.FromBigInt(new(big.Int).SetUint64(i), ctx.IntSort())
			case types.Int:
				i, err := strconv.ParseInt(v.Name, 10, intSize)
				if err != nil {
					panic(err)
				}
				return ctx.FromInt(i, ctx.IntSort())
			case types.Bool:
				b, err := strconv.ParseBool(v.Name)
				if err != nil {
					panic(err)
				}
				return ctx.FromBool(b)
			case types.Float64:
				f, err := strconv.ParseFloat(v.Name, floatSize)
				if err != nil {
					panic(err)
				}
				return ctx.FromFloat64(f, ctx.floatSort)
			case types.Complex128:
				c, err := strconv.ParseComplex(v.Name, complexSize)
				if err != nil {
					panic(err)
				}
				return ctx.FromComplex128(c)
			}
		}
		panic(fmt.Sprintf("unknown constant '%s' of type '%s'", v.Name, v.Type))
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
			panic(fmt.Sprintf("variable '%s' can't have different types ('%s' and '%s')", v.Name, oldV.Type, v.Type))
		}
		return
	}
	vars[v.Name] = v
}

func (bo BinOp) String() string {
	return fmt.Sprintf("%s == (%s %s %s)", bo.Result, bo.Left, bo.Op, bo.Right)
}

func (bo BinOp) Encode(ctx *EncodingContext) SymValue {
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
		}
	case "%":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Int).Eq(left.Mod(right.(z3.Int)))
		}
	case ">":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Bool).Eq(left.GT(right.(z3.Int)))
		case z3.Float:
			return res.(z3.Bool).Eq(left.GT(right.(z3.Float)))
		}
	case ">=":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Bool).Eq(left.GE(right.(z3.Int)))
		case z3.Float:
			return res.(z3.Bool).Eq(left.GE(right.(z3.Float)))
		}
	case "<":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Bool).Eq(left.LT(right.(z3.Int)))
		case z3.Float:
			return res.(z3.Bool).Eq(left.LT(right.(z3.Float)))
		}
	case "<=":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Bool).Eq(left.LE(right.(z3.Int)))
		case z3.Float:
			return res.(z3.Bool).Eq(left.LE(right.(z3.Float)))
		}
	case "==":
		switch left := left.(type) {
		case z3.Int:
			return res.(z3.Bool).Eq(left.Eq(right.(z3.Int)))
		case z3.Float:
			return res.(z3.Bool).Eq(left.IEEEEq(right.(z3.Float)))
		}
	case "<<":
		switch left := left.(type) {
		case z3.Int:
			leftBV := left.ToBV(intSize)
			rightBV := right.(z3.Int).ToBV(intSize)
			return res.(z3.Int).Eq(leftBV.Lsh(rightBV).SToInt())
		}
	case ">>":
		switch left := left.(type) {
		case z3.Int:
			leftBV := left.ToBV(intSize)
			rightBV := right.(z3.Int).ToBV(intSize)
			return res.(z3.Int).Eq(leftBV.SRsh(rightBV).SToInt())
		}
	case "^":
		switch left := left.(type) {
		case z3.Int:
			leftBV := left.ToBV(intSize)
			rightBV := right.(z3.Int).ToBV(intSize)
			return res.(z3.Int).Eq(leftBV.Xor(rightBV).SToInt())
		}
	case "&":
		switch left := left.(type) {
		case z3.Int:
			leftBV := left.ToBV(intSize)
			rightBV := right.(z3.Int).ToBV(intSize)
			return res.(z3.Int).Eq(leftBV.And(rightBV).SToInt())
		}
	case "|":
		switch left := left.(type) {
		case z3.Int:
			leftBV := left.ToBV(intSize)
			rightBV := right.(z3.Int).ToBV(intSize)
			return res.(z3.Int).Eq(leftBV.Or(rightBV).SToInt())
		}
	}
	panic(fmt.Sprintf("unknown binary operation '%s' for sort '%s'", bo.Op, left.Sort()))
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
	result := uo.Result.Encode(ctx)
	arg := uo.Arg.Encode(ctx)
	switch uo.Op {
	case "*":
		arg := arg.(*Pointer)
		switch result := result.(type) {
		case z3.Int:
			return result.Eq(ctx.valuesMemory[arg.t].Select(arg.addr).(z3.Int))
		case z3.Bool:
			return result.Eq(ctx.valuesMemory[arg.t].Select(arg.addr).(z3.Bool))
		case z3.Float:
			return result.Eq(ctx.valuesMemory[arg.t].Select(arg.addr).(z3.Float))
		}
	}
	panic(fmt.Sprintf("unknown unary operation '%s' for sort '%s'", uo.Op, arg.Sort()))
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
		}
		panic(fmt.Sprintf("unknown return sort '%s'", result.Sort()))
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
		arr := f.Args[0].Encode(ctx).(*SymArray)
		return f.Result.Encode(ctx).(z3.Int).Eq(ctx.arrayLenMemory[arr.t].Select(arr.addr).(z3.Int))
	}
	if function, ok := ctx.funcs[f.Name]; ok {
		var args []z3.Value
		for _, a := range f.Args {
			args = append(args, a.Encode(ctx).(z3.Value))
		}
		var res = f.Result.Encode(ctx)
		switch res := res.(type) {
		case z3.Int:
			return res.Eq(function.Apply(args...).(z3.Int))
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
	switch resT := c.Result.Type.(type) {
	case *types.Basic:
		switch resT.Kind() {
		case types.Int, types.Uint:
			switch argT := c.Arg.Type.(type) {
			case *types.Basic:
				switch argT.Kind() {
				case types.Int, types.Uint:
					return c.Result.Encode(ctx).(z3.Int).Eq(c.Arg.Encode(ctx).(z3.Int))
				}
			}
		case types.Float64:
			switch argT := c.Arg.Type.(type) {
			case *types.Basic:
				switch argT.Kind() {
				case types.Float64:
					return c.Result.Encode(ctx).(z3.Float).Eq(c.Arg.Encode(ctx).(z3.Float))
				case types.Int, types.Uint:
					return c.Result.Encode(ctx).(z3.Float).Eq(c.Arg.Encode(ctx).(z3.Int).ToBV(intSize).IEEEToFloat(ctx.floatSort))
				}
			}
		}
	}
	panic(fmt.Sprintf("unsupported conversion from '%s' to '%s'", c.Arg.Type, c.Result.Type))
}

func (c Convert) ScanVars(vars map[string]Var) {
	c.Arg.ScanVars(vars)
	c.Result.ScanVars(vars)
}

func (ia IndexAddr) String() string {
	return fmt.Sprintf("%s = &%s[%s]", ia.Result, ia.Array, ia.Index)
}

func (ia IndexAddr) Encode(ctx *EncodingContext) SymValue {
	res := ia.Result.Encode(ctx).(*Pointer).addr
	array := ia.Array.Encode(ctx).(*SymArray)
	index := ia.Index.Encode(ctx).(z3.Int)
	values := ctx.arrayValuesMemory[array.t].Select(array.addr).(z3.Array)
	value := values.Select(index).(z3.Uninterpreted)
	return res.Eq(value)
}

func (ia IndexAddr) ScanVars(vars map[string]Var) {
	ia.Result.ScanVars(vars)
	ia.Array.ScanVars(vars)
	ia.Index.ScanVars(vars)
}

func (fa FieldAddr) String() string {
	return fmt.Sprintf("%s = &%s#%d", fa.Result, fa.Object, fa.Field)
}

func (fa FieldAddr) Encode(ctx *EncodingContext) SymValue {
	panic("")
}

func (fa FieldAddr) ScanVars(vars map[string]Var) {
	fa.Result.ScanVars(vars)
	fa.Object.ScanVars(vars)
}

func toYaml(f Formula) string {
	d, err := yaml.Marshal(&f)
	if err != nil {
		panic(err)
	}
	return string(d)
}
