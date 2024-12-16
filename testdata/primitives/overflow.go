package main

type Overflow struct{}

func (o *Overflow) ShortOverflow(x int16, y int16) int16 {
	if y > 10 || y <= 0 {
		return 0
	}
	if x+y < 0 && x > 0 {
		panic("IllegalStateException")
	}
	return x + y
}
