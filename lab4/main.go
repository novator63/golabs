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

	ages := map[string]int{ //Задание 1
		"Кирилл":21,
		"Светлана": 18,
		"Владимир": 22,
		"Александр" : 30,
	}

	fmt.Println(ages) //Задание 1

	avg := AverageAge(ages) //Задание 2

	fmt.Printf("Средний возраст людей из списка: %.2f\n", avg) //Задание 2

	delete(ages, "Александр") //Задание 3

	for name, age := range ages{ //Задание 3
		fmt.Printf("%s: %d\n", name, age)
	}

	fmt.Println(ages)

	var word string

	fmt.Scanln(&word) //Задание 4
	fmt.Println(strings.ToUpper(word)) //Задание 4

	var a, b, c int
	fmt.Scan(&a, &b, &c) //Задание 5
	fmt.Println(a + b + c) //Задание 5

	var n int
	fmt.Print("Введите размерность массива: ") //Задание 6
	fmt.Scan(&n) //Задание 6

	nums := make([]int, n) //Задание 6
	fmt.Print("Задайте значения: ") //Задание 6
	for i := 0; i < n ; i++ { //Задание 6
		fmt.Scan(&nums[i]) //Задание 6
	}

	fmt.Println("В обратном порядке:") //Задание 6
	for i := n - 1; i >= 0; i-- { //Задание 6
		fmt.Print(nums[i], " ") //Задание 6
	}
}
