package main

import (
	"fmt"
	"sync"
)

func producer(n int, out chan <- int, wg *sync.WaitGroup){
	defer wg.Done()
	a, b := 0, 1
	for i := 0; i < n; i++{
		out <- a
		a,b = b, a+b
	}
	close(out)
}

func consumer(in <- chan int, wg *sync.WaitGroup){
	defer wg.Done()
	for x:=range in{
		fmt.Println(x)
	}
}

func main(){
	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan int)

	go producer(10, ch, &wg)
	go consumer(ch, &wg)

	wg.Wait()
}