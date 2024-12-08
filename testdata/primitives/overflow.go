package main

type Overflow struct{}

func (o *Overflow) ShortOverflow(x int16, y int) int {
	if y > 10 || y <= 0 {
		return 0
	}
	if int(x)+y < 0 && x > 0 {
		panic("IllegalStateException")
	}
	return int(x) + y
}
