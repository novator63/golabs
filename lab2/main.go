package main

import (
	"fmt"
)

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func average(n, a float64) float64 {
	return (n + a) / 2
}

func checkNumber(n int) string {
	if n > 0 {
		return "Positive"
	} else if n < 0 {
		return "Negative"
	}
	return "Zero"
}

func stringLength(str string) int {
	return len(str)
}

func main() {

	var input int

	fmt.Scan(&input)

	if input%2 == 0 {
		fmt.Printf("%d чет\n", input)
	} else {
		fmt.Println("нечет")
	}

	a := -4

	fmt.Printf("%s", checkNumber(a))

	for i := 0; i <= 10; i++ {
		fmt.Println(i)
	}

	fmt.Println(stringLength("Hello, world!"))

	rect := Rectangle{Width: 5, Height: 10}

	fmt.Println(rect.Area())

	fmt.Println(average(2, 4))
}
