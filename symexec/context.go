package symexec

import (
	"fmt"
	"go/types"
	"math"
	"math/big"

	"github.com/aclements/go-z3/z3"
)

const (
	intSize     = 64
	floatSize   = 64
	complexSize = 64

	resultSpecialVar = "$result"
)

type EncodingContext struct {
	*z3.Context

	asserts []z3.Bool

	vars     map[string]SymValue
	rawTypes map[string]z3.Sort

	varsUsed map[string]struct{}
	varCount map[string]int

	fieldsMemory      map[string][]z3.Array
	valuesMemory      map[string]z3.Array
	arrayValuesMemory map[string]z3.Array
	arrayLenMemory    map[string]z3.Array

	floatSort   z3.Sort
	complexSort z3.Sort
	stringSort  z3.Sort

	addrSort z3.Sort
}

type NamedStruct struct {
	*types.Struct
	Name string
}

func (ctx *EncodingContext) AddType(t types.Type) z3.Sort {
	if _, ok := ctx.rawTypes[t.String()]; !ok {
		switch t := t.(type) {
		case *types.Basic:
			switch t.Kind() {
			case types.Int:
				ctx.rawTypes[t.String()] = ctx.IntSort()
			case types.Bool:
				ctx.rawTypes[t.String()] = ctx.BoolSort()
			case types.Float64:
				ctx.rawTypes[t.String()] = ctx.floatSort
			case types.Complex128:
				ctx.rawTypes[t.String()] = ctx.addrSort // TODO: complex number representation as z3.Sort
			case types.String:
				ctx.rawTypes[t.String()] = ctx.addrSort // TODO: string representation as z3.Sort
			default:
				panic(fmt.Sprintf("unknown basic type '%s'", t))
			}
		case *types.Pointer:
			elemT := ctx.AddType(t.Elem())
			ctx.valuesMemory[t.String()] = ctx.Const(fmt.Sprintf("$<%s>Memory", t), ctx.ArraySort(ctx.addrSort, elemT)).(z3.Array)
			ctx.rawTypes[t.String()] = ctx.addrSort
		case *types.Slice:
			elemT := ctx.AddType(types.NewPointer(t.Elem()))
			ctx.arrayValuesMemory[t.String()] = ctx.Const(
				fmt.Sprintf("$<%s>ValuesMemory", t),
				ctx.ArraySort(ctx.addrSort, ctx.ArraySort(ctx.IntSort(), elemT)),
			).(z3.Array)
			ctx.arrayLenMemory[t.String()] = ctx.Const(
				fmt.Sprintf("$<%s>LenMemory", t),
				ctx.ArraySort(ctx.addrSort, ctx.IntSort()),
			).(z3.Array)
			ctx.rawTypes[t.String()] = ctx.addrSort
		case *types.Struct:
			var fields []z3.Array
			for i := 0; i < t.NumFields(); i++ {
				f := t.Field(i)
				elemT := ctx.AddType(types.NewPointer(f.Type()))
				fieldArray := ctx.Const(
					fmt.Sprintf("$<%s.%s:%s>Memory", t, f.Name(), f.Type()),
					ctx.ArraySort(ctx.addrSort, elemT),
				).(z3.Array)
				fields = append(fields, fieldArray)
			}
			ctx.fieldsMemory[t.String()] = fields
			ctx.rawTypes[t.String()] = ctx.addrSort
		case NamedStruct:
			var fields []z3.Array
			for i := 0; i < t.NumFields(); i++ {
				f := t.Field(i)
				elemT := ctx.AddType(types.NewPointer(f.Type()))
				fieldArray := ctx.Const(
					fmt.Sprintf("$<%s.%s:%s>Memory", t.Name, f.Name(), f.Type()),
					ctx.ArraySort(ctx.addrSort, elemT),
				).(z3.Array)
				fields = append(fields, fieldArray)
			}
			ctx.fieldsMemory[t.Name] = fields
			ctx.rawTypes[t.Name] = ctx.addrSort
		case *types.Named:
			ctx.AddType(NamedStruct{Struct: t.Underlying().(*types.Struct), Name: t.String()})
		default:
			panic(fmt.Sprintf("unknown type '%s'", t))
		}
	}
	return ctx.rawTypes[t.String()]
}

func (ctx *EncodingContext) AddVar(name string, z3name string, t types.Type) {
	switch t := t.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Int:
			i := ctx.IntConst(z3name)
			ctx.vars[name] = i
			ctx.asserts = append(ctx.asserts, i.LE(ctx.FromInt(math.MaxInt, ctx.IntSort()).(z3.Int)))
			ctx.asserts = append(ctx.asserts, i.GE(ctx.FromInt(math.MinInt, ctx.IntSort()).(z3.Int)))
		case types.Uint:
			i := ctx.IntConst(z3name)
			ctx.vars[name] = i
			ctx.asserts = append(ctx.asserts, i.LE(ctx.FromBigInt(new(big.Int).SetUint64(math.MaxUint), ctx.IntSort()).(z3.Int)))
			ctx.asserts = append(ctx.asserts, i.GE(ctx.FromInt(0, ctx.IntSort()).(z3.Int)))
		case types.Bool:
			ctx.vars[name] = ctx.BoolConst(z3name)
		case types.Float64:
			ctx.vars[name] = ctx.Const(z3name, ctx.floatSort)
		case types.Complex128:
			ctx.vars[name] = ctx.ComplexConst(z3name)
		case types.String:
			ctx.vars[name] = ctx.StringConst(z3name)
		}
	case *types.Pointer:
		ctx.vars[name] = ctx.PointerConst(z3name, t.String(), t.Elem().String())
	case *types.Slice:
		ctx.vars[name] = ctx.SymArrayConst(z3name, t.String())
	case *types.Struct:
		ctx.vars[name] = ctx.SymStructConst(z3name, t.String())
	default:
		panic(fmt.Sprintf("variable '%s' of unknown type '%s'", name, t))
	}
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

func (ctx *EncodingContext) SymArrayConst(name string, t string) *SymArray {
	return &SymArray{
		addr: ctx.Const(name, ctx.addrSort).(z3.Uninterpreted),
		t:    t,
		sort: ctx.rawTypes[t],
	}
}

func (ctx *EncodingContext) PointerConst(name string, t string, elem string) *Pointer {
	return &Pointer{
		addr: ctx.Const(name, ctx.addrSort).(z3.Uninterpreted),
		t:    t,
		elem: elem,
		sort: ctx.rawTypes[t],
	}
}

func (ctx *EncodingContext) SymStructConst(name string, t string) *SymStruct {
	return &SymStruct{
		addr: ctx.Const(name, ctx.addrSort).(z3.Uninterpreted),
		t:    t,
		sort: ctx.rawTypes[t],
	}
}

func (ctx *EncodingContext) FromComplex128(c complex128) *Complex {
	return &Complex{
		real: ctx.FromFloat64(real(c), ctx.floatSort),
		imag: ctx.FromFloat64(imag(c), ctx.floatSort),
		sort: ctx.complexSort,
	}
}
