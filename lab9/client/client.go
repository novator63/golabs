package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

const baseURL = "http://localhost:8080"

var userTokens = make(map[string]string) // username → token
var currentUser string                   // имя текущего активного пользователя

func main() {
	for {
		fmt.Println("\n=== КЛИЕНТ REST API ===")
		fmt.Println("1. Авторизация (логин)")
		fmt.Println("2. Сменить активного пользователя")
		fmt.Println("3. Показать всех пользователей")
		fmt.Println("4. Найти пользователя по ID")
		fmt.Println("5. Добавить нового пользователя")
		fmt.Println("6. Обновить данные пользователя")
		fmt.Println("7. Удалить пользователя")
		fmt.Println("8. Показать сохранённые токены")
		fmt.Println("0. Выход")

		fmt.Print("Выберите действие: ")
		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			login()
		case 2:
			switchUser()
		case 3:
			getAllUsers()
		case 4:
			getUserByID()
		case 5:
			createUser()
		case 6:
			updateUser()
		case 7:
			deleteUser()
		case 8:
			printTokens()
		case 0:
			fmt.Println("Выход.")
			os.Exit(0)
		default:
			fmt.Println("Неверный выбор.")
		}
	}
}

// ---------------- Авторизация ----------------
func login() {
	var username, password string
	fmt.Print("Введите логин: ")
	fmt.Scan(&username)
	fmt.Print("Введите пароль: ")
	fmt.Scan(&password)

	data, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})

	resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Println("Ошибка авторизации:", string(body))
		return
	}

	var result map[string]string
	json.Unmarshal(body, &result)
	token := result["token"]

	userTokens[username] = token
	currentUser = username

	fmt.Printf("✅ Пользователь '%s' авторизован.\n", username)
	fmt.Printf("Токен сохранён: %s\n", token)
}

// ---------------- Смена пользователя ----------------
func switchUser() {
	if len(userTokens) == 0 {
		fmt.Println("Нет авторизованных пользователей.")
		return
	}

	fmt.Println("Доступные пользователи:")
	for u := range userTokens {
		fmt.Println("-", u)
	}
	fmt.Print("Введите имя пользователя для активации: ")
	var name string
	fmt.Scan(&name)

	if _, ok := userTokens[name]; !ok {
		fmt.Println("Такого пользователя нет в списке.")
		return
	}

	currentUser = name
	fmt.Printf("🔄 Активный пользователь: %s\n", currentUser)
}

// ---------------- CRUD ----------------
func getAllUsers() {
	req, _ := http.NewRequest("GET", baseURL+"/users", nil)
	addAuth(req)
	send(req)
}

func getUserByID() {
	var id int
	fmt.Print("Введите ID: ")
	fmt.Scan(&id)
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/%d", baseURL, id), nil)
	addAuth(req)
	send(req)
}

func createUser() {
	var u User
	fmt.Print("Имя: ")
	fmt.Scan(&u.Name)
	fmt.Print("Возраст: ")
	fmt.Scan(&u.Age)
	fmt.Print("Логин: ")
	fmt.Scan(&u.Username)
	fmt.Print("Пароль: ")
	fmt.Scan(&u.Password)

	data, _ := json.Marshal(u)
	req, _ := http.NewRequest("POST", baseURL+"/users", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	addAuth(req)
	send(req)
}

func updateUser() {
	var id int
	var u User
	fmt.Print("Введите ID пользователя для обновления: ")
	fmt.Scan(&id)
	fmt.Print("Новое имя: ")
	fmt.Scan(&u.Name)
	fmt.Print("Новый возраст: ")
	fmt.Scan(&u.Age)

	data, _ := json.Marshal(u)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/users/%d", baseURL, id), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	addAuth(req)
	send(req)
}

func deleteUser() {
	var id int
	fmt.Print("Введите ID для удаления: ")
	fmt.Scan(&id)
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%d", baseURL, id), nil)
	addAuth(req)
	send(req)
}

// ---------------- Вспомогательные функции ----------------
func addAuth(req *http.Request) {
	if currentUser == "" {
		fmt.Println("⚠️ Сначала выполните авторизацию.")
		return
	}
	token := userTokens[currentUser]
	req.Header.Set("Authorization", "Bearer "+token)
}

func send(req *http.Request) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("\nОтвет (%d): %s\n", resp.StatusCode, string(body))
}

func printTokens() {
	if len(userTokens) == 0 {
		fmt.Println("Нет сохранённых токенов.")
		return
	}
	fmt.Println("\nСохранённые токены:")
	for user, token := range userTokens {
		active := ""
		if user == currentUser {
			active = "(активен)"
		}
		fmt.Printf("- %s: %s %s\n", user, token, active)
	}
}
