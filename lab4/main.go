package main

import (
	"fmt"
	"strings"
)

func AverageAge(people map[string] int) float64{
	sum := 0

	for _, age := range people{
		sum += age
	}

	return float64(sum) / float64(len(people))
}

func main(){

	ages := map[string]int{
		"Кирилл":21,
		"Светлана": 18,
		"Владимир": 22,
		"Александр" : 30,
	}

	fmt.Println(ages)

	avg := AverageAge(ages)

	fmt.Printf("Средний возраст людей из списка: %.2f\n", avg)

	delete(ages, "Александр")

	for name, age := range ages{
		fmt.Printf("%s: %d\n", name, age)
	}

	var word string

	fmt.Scanln(&word)
	fmt.Println(strings.ToUpper(word))

	var a, b, c int
	fmt.Scan(&a, &b, &c)
	fmt.Println(a + b + c)

	var n int
	fmt.Print("Введите размерность массива: ")
	fmt.Scan(&n)

	nums := make([]int, n)
	fmt.Print("Задайте значения: ")
	for i := 0; i < n ; i++ {
		fmt.Scan(&nums[i])
	}

	fmt.Println("В обратном порядке:")
	for i := n - 1; i >= 0; i-- {
		fmt.Print(nums[i], " ")
	}
}