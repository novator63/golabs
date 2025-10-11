package main

import (
	"fmt"
	// "lab3/mathutils"
	// "lab3/stringutils"
)

func main() {

	//Задание 1-3
	// n := 0
	// fmt.Scanf("%d", &n)
	// fmt.Printf("%d! = %d\n", n, mathutils.Factorial(n))

	// fmt.Println(stringutils.Reverse("Hello, world!"))

	//Задиние 4
	// var arr [5]int

	// for i := 0; i < len(arr); i++ {
	// 	arr[i] = i
	// }

	// fmt.Println(arr)

	//Задание 5

	// slice := []int{1, 2, 3, 4, 5}

	// slice = append(slice, 6, 7, 8)
	// fmt.Println("После добавления 6,7,8:", slice)

	// slice = slice[1:6]
	// fmt.Println("После обрезки [1:6]:", slice)

	// m := 2

	// slice = append(slice[:m], slice[m+1:]...)
	// fmt.Println("После удаления элемента по индексу", m, ":", slice)

	// value := 99
	// slice = append(slice[:m], append([]int{value}, slice[m:]...)...)
	// fmt.Println("После добавления элемента 99 по индексу", m, ":", slice)

	//Задиние 6
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
