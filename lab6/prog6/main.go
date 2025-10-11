package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

// Task — входные данные (строка)
type Task struct {
	Line string
}

// Result — результат обработки
type Result struct {
	Original string
	Reversed string
}

// reverseString — реверс строки (руны, чтобы работали юникод-символы)
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// worker — воркер читает задачи, пишет результаты
func worker(id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		rev := reverseString(task.Line)
		results <- Result{Original: task.Line, Reversed: rev}
	}
}

func main() {
	// === Входные данные ===
	inputFile := "input.txt"   // файл со строками (по одной на строку)
	outputFile := "output.txt" // файл для результата

	// Сколько воркеров запустить — задаёт пользователь
	var numWorkers int
	fmt.Print("Введите количество воркеров: ")
	fmt.Scan(&numWorkers)

	// === Каналы ===
	tasks := make(chan Task)
	results := make(chan Result)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	// Запускаем пул воркеров
	for i := 0; i < numWorkers; i++ {
		go worker(i, tasks, results, &wg)
	}

	// Читаем строки из файла и отправляем в канал задач
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

	// Горутинa для закрытия results после завершения всех workers
	go func() {
		wg.Wait()
		close(results)
	}()

	// === Вывод результатов ===
	out, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer out.Close()

	for r := range results {
		line := fmt.Sprintf("Исходная: %s | Реверс: %s\n", r.Original, r.Reversed)
		fmt.Print(line)       // вывод в консоль
		out.WriteString(line) // запись в файл
	}
}

/*
Принцип работы программы (пул воркеров с реверсом строк):

1. Пользователь задаёт количество воркеров (горутин).
   Эти воркеры будут параллельно обрабатывать задачи.

2. Каналы:
   - tasks — канал для входных задач (строки из файла).
   - results — канал для результатов (оригинал + реверс строки).

3. Каждый воркер берёт строку из tasks, делает её реверс,
   и кладёт результат в results.

4. Отдельная горутина читает строки из input.txt построчно
   и отправляет их в канал tasks. После окончания чтения
   tasks закрывается.

5. После завершения всех воркеров (через WaitGroup)
   канал results закрывается.

6. Главный цикл программы читает из results и:
   - выводит результат в консоль,
   - одновременно записывает результат в output.txt.

Итог: строки обрабатываются параллельно несколькими воркерами,
каждая строка реверсируется, а результаты сохраняются и печатаются.
*/
