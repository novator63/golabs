package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	numbers := make(chan int)      
	toParity := make(chan int)
	messages := make(chan string)

	go func() {
		for {
			numbers <- rand.Intn(100)
			time.Sleep(150 * time.Millisecond)
		}
	}()

	go func() {
		for n := range toParity {
			if n%2 == 0 {
				messages <- fmt.Sprintf("Число %d — чётное", n)
			} else {
				messages <- fmt.Sprintf("Число %d — нечётное", n)
			}
		}
	}()

	for {
		select {
		case n := <-numbers:
			fmt.Printf("Сгенерированное число: %d\n", n)
			toParity <- n

		case msg := <-messages:
			fmt.Println(msg)

		case <-time.After(1 * time.Second):
			fmt.Println("Таймаут ожидания...")
		}
	}
}
