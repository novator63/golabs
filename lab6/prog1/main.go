package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func factorial(n int) int {
	res := 1
	for i := 2; i <= n; i++ {
		time.Sleep(200 * time.Microsecond)
		res *= i
		fmt.Printf("[factorial] умножил на %d, результат = %d\n", i, res)
	}
	return res
}

func generateRandom(count int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	out := make([]int, 0, count)
	for i := 0; i < count; i++ {
		time.Sleep(120 * time.Millisecond)
		n := r.Intn(100)
		out = append(out, n)
		fmt.Printf("[random] сгенерировали %d\n", n)
	}
	return out
}

func sumSeries(n int) int{
	sum := 0
	for i := 0; i < n; i++{
		time.Sleep(100 * time.Microsecond)
		sum += i
		fmt.Printf("[series] добавил %d, текущая сумма %d\n", i, sum)
	}
	return sum
}

func main() {

	var (
		wg sync.WaitGroup
		facRes int
		rndRes []int
		sumRes int
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		fmt.Println("[factorial] старт")
		facRes = factorial(5)
		fmt.Println("[factorial] готово")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("[random] старт")
		rndRes = generateRandom(8)
		fmt.Println("[random] готово")
	}()

		go func() {
		defer wg.Done()
		fmt.Println("[series] старт")
		sumRes = sumSeries(12)
		fmt.Println("[series] готово")
	}()

	wg.Wait()
	fmt.Println("Все горутины завершены:")
	fmt.Printf("ИТОГ факториала: %d\n", facRes)
	fmt.Printf("ИТОГ случайных чисел: %v\n", rndRes)
	fmt.Printf("ИТОГ суммы ряда: %d\n", sumRes)
}
