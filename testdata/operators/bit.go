package main

type BitOperators struct{}

func (b *BitOperators) Complement(x int) bool {
	return ^x == 1
}

func (b *BitOperators) Xor(x, y int) bool {
	return (x ^ y) == 0
}

func (b *BitOperators) Or(val int) bool {
	return (val | 7) == 15
}

func (b *BitOperators) And(value int) bool {
	return (value & (value - 1)) == 0
}

func (b *BitOperators) BooleanNot(boolA, boolB bool) int {
	d := boolA && boolB
	e := (!boolA) || boolB // Используем корректные имена переменных
	if d && e {
		return 100
	}
	return 200
}

func (b *BitOperators) BooleanXor(valA, valB bool) bool {
	return valA != valB
}

func (b *BitOperators) BooleanOr(aVal, bVal bool) bool {
	return aVal || bVal
}

func (b *BitOperators) BooleanAnd(aVal, bVal bool) bool {
	return aVal && bVal
}

func (b *BitOperators) BooleanXorCompare(aBool, bBool bool) int {
	if aBool != bBool {
		return 1
	}
	return 0
}

func (b *BitOperators) Shl(num int) bool {
	return (num << 1) == 2
}

func (b *BitOperators) ShlLong(longNum int64) bool {
	return (longNum << 1) == 2
}

func (b *BitOperators) ShlWithBigLongShift(shift int64) int {
	if shift < 40 {
		return 1
	}
	if (0x77777777 << shift) == 0x77777770 {
		return 2
	}
	return 3
}

func (b *BitOperators) Shr(num int) bool {
	return (num >> 1) == 1
}

func (b *BitOperators) ShrLong(longVal int64) bool {
	return (longVal >> 1) == 1
}

func (b *BitOperators) Ushr(unsignedVal int) bool {
	return (unsignedVal >> 1) == 1 // В Go сдвиг вправо всегда беззнаковый
}

func (b *BitOperators) UshrLong(unsignedLongVal int64) bool {
	return (unsignedLongVal >> 1) == 1 // В Go сдвиг вправо всегда беззнаковый
}

func (b *BitOperators) Sign(num int) int {
	i := (num >> 31) | (-num >> 31)

	if i > 0 {
		return 1
	} else if i == 0 {
		return 0
	} else {
		return -1
	}
}
