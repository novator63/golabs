package main

import (
	"fmt"
	"time"
)

func sumAndDiff(a, b float64) (float64, float64) {
	sum := a + b
	diff := a - b
	return sum, diff
}

func main() {
	now := time.Now()
	fmt.Println(now)

	var a int = 10
	var k float64 = 3.13
	var s string = "Hello"
	var b bool = true

	fmt.Println(a, k, s, b)

	x := 100
	y := 2.71
	msg := "World"
	z := false

	fmt.Println(x, y, msg, z)

	m, n := 15, 4

	fmt.Printf("%d + %d = %d\n", m, n, m+n)
	fmt.Printf("%d - %d = %d\n", m, n, m-n)
	fmt.Printf("%d * %d = %d\n", m, n, m*n)
	fmt.Printf("%d / %d = %d\n", m, n, m/n)
	fmt.Printf("%d %% %d = %d\n", m, n, m%n)

	sum, diff := sumAndDiff(7.2, 2.8)
	fmt.Println(sum, diff)

	x1, x2, x3 := 10.0, 20.0, 30.0
	avg := (x1 + x2 + x3) / 3
	fmt.Println(avg)

}
