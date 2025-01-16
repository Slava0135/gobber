# Gobber

Go (dynamic/static) symbolic execution using Z3 SMT solver.

>Not (even close to) production ready. Educational purpose only!

## Does it work?

Kinda.

There are tests [here](./symexec/symexec_test.go). Some of them work. Mostly ported from <https://github.com/UnitTestBot/usvm> tests.

For example:

```go
package main

func pushPopIncrementality(j int) int {
    result := j

    for i := 1; i <= 10; i++ {
        result += i
    }

    if result%2 == 0 {
        result++
    }

    return result
}
```

...will produce following output:

```text
=== RUN   TestDynamic_PushPop
:: building SSA graph for file 'push_pop.go'
:: analyzing function 'pushPopIncrementality'
:: printing SSA blocks
0 ->
  [      JUMP] jump 1
1 ->
  [       PHI] t0:int <-- phi [0: j, 2: t3] #result
  [       PHI] t1:int <-- phi [0: 1:int, 2: t4] #i
  [     BINOP] t2:bool <-- t1 <= 10:int
  [        IF] if t2 goto 2 else 3
2 ->
  [     BINOP] t3:int <-- t0 + t1
  [     BINOP] t4:int <-- t1 + 1:int
  [      JUMP] jump 1
3 ->
  [     BINOP] t5:int <-- t0 % 2:int
  [     BINOP] t6:bool <-- t5 == 0:int
  [        IF] if t6 goto 4 else 5
4 ->
  [     BINOP] t7:int <-- t0 + 1:int
  [      JUMP] jump 5
5 ->
  [       PHI] t8:int <-- phi [3: t0, 4: t7] #result
  [    RETURN] return t8
:: execute
found solution for path: [0 1 2 1 2 1 2 1 2 1 2 1 2 1 2 1 2 1 2 1 2 1 3 4 5]
t4~7 -> 9
t2~10 -> false
t2~7 -> true
t2~2 -> true
t1~7 -> 8
t4~9 -> 11
t6 -> true
$result -> 1
...
```

...and generate tests in same directory, covering all paths(!):

```go
package main

import (
	"math"
	"testing"
)

var (
	_ = testing.Main
	_ = math.Abs
)

func Test_pushPopIncrementality_1(t *testing.T) {
	j := -55
	got := pushPopIncrementality(j)
	want := 1
	if got != want {
		t.Errorf("pushPopIncrementality(j) = %v; want %v", got, want)
	}
}

func Test_pushPopIncrementality_2(t *testing.T) {
	j := -54
	got := pushPopIncrementality(j)
	want := 1
	if got != want {
		t.Errorf("pushPopIncrementality(j) = %v; want %v", got, want)
	}
}
```

You can run these as normal tests, to verify they are correct (and that they cover all lines).

Another example:

```go
func mixedOperations(a int, b float64) float64 {
	var result float64

	if a%2 == 0 {
		result = float64(a) + b
	} else {
		result = float64(a) - b
	}

	if result < 10 {
		result *= 2
	} else {
		result /= 2
	}

	return result
}
```

```go
package main

import (
	"math"
	"testing"
)

var (
	_ = testing.Main
	_ = math.Abs
)

func Test_mixedOperations_1(t *testing.T) {
	a := 74907589934418923
	b_bits := uint64(0x4518081bff960000) // 7263128770402943422693376.000000
	b := math.Float64frombits(b_bits)
	got := mixedOperations(a, b)
	want_bits := uint64(0xc528081bfb6d7fc6) // -14526257390990705937088512.000000
	want := math.Float64frombits(want_bits)
	if math.Abs(got - want) > 1e-6 && !(math.IsNaN(got) && math.IsNaN(want)) {
		t.Errorf("mixedOperations(a, b) = %v; want %v", got, want)
	}
}

func Test_mixedOperations_2(t *testing.T) {
	a := 704374636545
	b_bits := uint64(0xe030200000000000) // some very big negative number
	b := math.Float64frombits(b_bits)
	got := mixedOperations(a, b)
	want_bits := uint64(0x6020200000000000) // some very big positive number
	want := math.Float64frombits(want_bits)
	if math.Abs(got - want) > 1e-6 && !(math.IsNaN(got) && math.IsNaN(want)) {
		t.Errorf("mixedOperations(a, b) = %v; want %v", got, want)
	}
}

func Test_mixedOperations_3(t *testing.T) {
	a := -11217466889863938
	b_bits := uint64(0x434c4224800000a4) // 15908047763276104.000000
	b := math.Float64frombits(b_bits)
	got := mixedOperations(a, b)
	want_bits := uint64(0x4320aa0ef6bffe46) // 2345290436706083.000000
	want := math.Float64frombits(want_bits)
	if math.Abs(got - want) > 1e-6 && !(math.IsNaN(got) && math.IsNaN(want)) {
		t.Errorf("mixedOperations(a, b) = %v; want %v", got, want)
	}
}

func Test_mixedOperations_4(t *testing.T) {
	a := 0
	b_bits := uint64(0x5b00804e002c) // 0.000000
	b := math.Float64frombits(b_bits)
	got := mixedOperations(a, b)
	want_bits := uint64(0xb601009c0058) // 0.000000
	want := math.Float64frombits(want_bits)
	if math.Abs(got - want) > 1e-6 && !(math.IsNaN(got) && math.IsNaN(want)) {
		t.Errorf("mixedOperations(a, b) = %v; want %v", got, want)
	}
}
```

## Install

```sh
git clone --recurse-submodules
```

Install z3 as explained in go-z3 package. You can throw away fork in this repository (if you don't need unsat cores), just remove some code in `./constraints` that depends on it.
