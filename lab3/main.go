package main

import (
	"fmt"
	"lab3/mathutils"
	"lab3/stringutils"
)

func main() {

	n := 0
	fmt.Scanf("%d", &n)
	fmt.Printf("%d! = %d\n", n, mathutils.Factorial(n))

	fmt.Println(stringutils.Reverse("Hello, world!"))

	var arr [5]int

	for i := 0; i < len(arr); i++ {
		arr[i] = i
	}

	fmt.Println(arr)

	slice := []int{1, 2, 3, 4, 5}

	slice = append(slice, 6, 7, 8)

	fmt.Println(slice)

	slice = slice[1:6]

	fmt.Println(slice)

	m := 2

	slice = append(slice[:m], slice[m+1:]...)

	fmt.Println(slice)

	notes := []string{"каждый", "охотник", "желает", "знать", "где", "сидит", "фазан"}

	var longest string
	maxLen := 0

	for _, word := range notes {
		if len([]rune(word)) > maxLen {
			maxLen = len([]rune(word))
			longest = word
		}
	}

	fmt.Println(longest)
}
