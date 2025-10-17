package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Структура пользователя
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

const baseURL = "http://localhost:8080/users"

// ---------------------------
// Основное меню
// ---------------------------
func main() {
	for {
		fmt.Println("\n=== КЛИЕНТ REST API ===")
		fmt.Println("1. Показать всех пользователей")
		fmt.Println("2. Найти пользователя по ID")
		fmt.Println("3. Добавить нового пользователя")
		fmt.Println("4. Обновить данные пользователя")
		fmt.Println("5. Удалить пользователя")
		fmt.Println("0. Выход")

		fmt.Print("Выберите действие: ")
		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			getAllUsers()
		case 2:
			getUserByID()
		case 3:
			createUser()
		case 4:
			updateUser()
		case 5:
			deleteUser()
		case 0:
			fmt.Println("Выход из программы.")
			os.Exit(0)
		default:
			fmt.Println("Неверный выбор, попробуйте снова.")
		}
	}
}

// ---------------------------
// 1. Получение всех пользователей
// ---------------------------
func getAllUsers() {
	resp, err := http.Get(baseURL)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("\nОтвет сервера:\n", string(body))
}

// ---------------------------
// 2. Получение пользователя по ID
// ---------------------------
func getUserByID() {
	var id int
	fmt.Print("Введите ID пользователя: ")
	fmt.Scan(&id)

	resp, err := http.Get(fmt.Sprintf("%s/%d", baseURL, id))
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("\nОтвет сервера:\n", string(body))
}

// ---------------------------
// 3. Добавление пользователя
// ---------------------------
func createUser() {
	var user User
	fmt.Print("Введите имя: ")
	fmt.Scan(&user.Name)
	fmt.Print("Введите возраст: ")
	fmt.Scan(&user.Age)

	data, _ := json.Marshal(user)
	resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("\nОтвет сервера:\n", string(body))
}

// ---------------------------
// 4. Обновление пользователя
// ---------------------------
func updateUser() {
	var id int
	var user User
	fmt.Print("Введите ID пользователя для обновления: ")
	fmt.Scan(&id)
	fmt.Print("Введите новое имя: ")
	fmt.Scan(&user.Name)
	fmt.Print("Введите новый возраст: ")
	fmt.Scan(&user.Age)

	data, _ := json.Marshal(user)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/%d", baseURL, id), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("\nОтвет сервера:\n", string(body))
}

// ---------------------------
// 5. Удаление пользователя
// ---------------------------
func deleteUser() {
	var id int
	fmt.Print("Введите ID пользователя для удаления: ")
	fmt.Scan(&id)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%d", baseURL, id), nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Пользователь успешно удалён.")
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Ответ сервера:", string(body))
	}
}
