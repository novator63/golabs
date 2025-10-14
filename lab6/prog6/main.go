package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

type Task struct {
	Line string
}

type Result struct {
	Original string
	Reversed string
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func worker(id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		rev := reverseString(task.Line)
		results <- Result{Original: task.Line, Reversed: rev}
	}
}

func main() {
	inputFile := "input.txt" 
	outputFile := "output.txt"

	var numWorkers int
	fmt.Print("Введите количество воркеров: ")
	fmt.Scan(&numWorkers)

	tasks := make(chan Task)
	results := make(chan Result)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go worker(i, tasks, results, &wg)
	}

	go func() {
		file, err := os.Open(inputFile)
		if err != nil {
			fmt.Println("Ошибка открытия файла:", err)
			close(tasks)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			tasks <- Task{Line: scanner.Text()}
		}
		close(tasks)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer out.Close()

	for r := range results {
		line := fmt.Sprintf("Исходная: %s | Реверс: %s\n", r.Original, r.Reversed)
		fmt.Print(line) 
		out.WriteString(line)
	}
}