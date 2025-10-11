package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	serverAddr := "127.0.0.1:9000"
	conn, err := net.Dial("tcp", serverAddr) //подключение к серверу
	if err != nil {
		log.Println("Ошибка подключения:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Подключено к серверу", serverAddr)
	fmt.Println("Введите сообщение (или 'exit' для выхода):")

	consoleReader := bufio.NewReader(os.Stdin) //создаёт сканер для чтения с консоли

	for {
		fmt.Print("> ")
		text, _ := consoleReader.ReadString('\n')
		text = strings.TrimSpace(text) // убираем \n и пробелы

		if text == "exit" {
			log.Println("Выход из клиента.")
			break
		}

		//отправляет сообщение серверу
		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			log.Println("Ошибка при отправке:", err)
			break
		}

		//читает ответ от сервера
		reply := make([]byte, 1024)
		n, err := conn.Read(reply)
		if err != nil {
			log.Println("Ошибка при чтении ответа:", err)
			break
		}

		fmt.Printf("Ответ сервера: %s\n", string(reply[:n]))
	}
}
