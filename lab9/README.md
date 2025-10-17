REST API Клиент–Серверное Приложение на Go
Общее описание

Данное приложение представляет собой полноценную реализацию клиент–серверной архитектуры на языке Go.
Серверная часть реализует REST API для управления пользователями (CRUD операции) и авторизацию с использованием токенов доступа.
Клиентская часть является консольной программой, которая взаимодействует с сервером через HTTP-запросы, выполняя операции с пользователями и авторизацию.

Приложение использует базу данных PostgreSQL для хранения данных о пользователях и активных сессиях.

1. Серверная часть
1.1 Назначение

Сервер предоставляет REST API, который позволяет:

выполнять авторизацию пользователей (/login);

создавать, читать, изменять и удалять записи пользователей (/users, /users/{id});

обеспечивать доступ к защищённым маршрутам только при наличии действительного токена.

1.2 Общая структура программы
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)


Используемые библиотеки:

database/sql — стандартный интерфейс для работы с SQL-базами данных;

github.com/lib/pq — драйвер PostgreSQL;

github.com/gorilla/mux — маршрутизатор для управления HTTP-запросами;

github.com/google/uuid — генерация уникальных идентификаторов (токенов);

encoding/json — кодирование и декодирование данных в формате JSON.

1.3 Структура данных пользователя
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
}


ID — идентификатор пользователя в таблице;

Username и Password — учетные данные для входа;

Name и Age — дополнительные сведения о пользователе;

тег omitempty исключает поле Password из JSON-ответа для защиты данных.

1.4 Подключение к базе данных

Сервер подключается к PostgreSQL при запуске:

connStr := "user=postgres password=postgres dbname=usersdb sslmode=disable"
db, err = sql.Open("postgres", connStr)
if err != nil {
	log.Fatal(err)
}


Далее выполняется проверка соединения методом db.Ping().

1.5 Маршруты

После подключения создаётся маршрутизатор:

r := mux.NewRouter()
r.HandleFunc("/login", loginHandler).Methods("POST")
api := r.PathPrefix("/").Subrouter()
api.Use(authMiddleware)
api.HandleFunc("/users", getUsers).Methods("GET")
api.HandleFunc("/users/{id}", getUser).Methods("GET")
api.HandleFunc("/users", createUser).Methods("POST")
api.HandleFunc("/users/{id}", updateUser).Methods("PUT")
api.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")


/login — маршрут для авторизации (открытый);

остальные маршруты /users защищены middleware authMiddleware.

1.6 Авторизация и создание токена

Функция loginHandler обрабатывает вход пользователя:

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)

	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username=$1 AND password=$2",
		creds.Username, creds.Password).Scan(&userID)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusUnauthorized, "Неверные учетные данные")
		return
	}

	token := uuid.New().String()
	db.Exec("INSERT INTO sessions (user_id, token, created_at) VALUES ($1, $2, $3)",
		userID, token, time.Now())

	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}


Алгоритм:

Сервер принимает логин и пароль.

Проверяет пользователя в таблице users.

Если учётные данные корректны — генерируется уникальный token (UUID).

Токен сохраняется в таблицу sessions.

Клиенту возвращается токен в JSON-ответе.

1.7 Проверка токена (middleware)

Каждый запрос к защищённым маршрутам проходит проверку:

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			respondWithError(w, http.StatusUnauthorized, "Отсутствует токен")
			return
		}

		var userID int
		err := db.QueryRow("SELECT user_id FROM sessions WHERE token=$1", token).Scan(&userID)
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusUnauthorized, "Недействительный токен")
			return
		}

		next.ServeHTTP(w, r)
	})
}


Порядок работы:

Извлекается заголовок Authorization.

Токен проверяется в таблице sessions.

Если токен найден — запрос передаётся дальше к целевому обработчику.

Если токен отсутствует или неверен — сервер возвращает ошибку 401 Unauthorized.

1.8 CRUD операции
Получение всех пользователей
func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT id, name, age, username FROM users ORDER BY id")
	defer rows.Close()
	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Age, &u.Username)
		users = append(users, u)
	}
	respondWithJSON(w, http.StatusOK, users)
}

Получение пользователя по ID
func getUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var u User
	db.QueryRow("SELECT id, name, age, username FROM users WHERE id=$1", id).Scan(&u.ID, &u.Name, &u.Age, &u.Username)
	respondWithJSON(w, http.StatusOK, u)
}

Создание пользователя
func createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	db.Exec("INSERT INTO users (name, age, username, password) VALUES ($1, $2, $3, $4)",
		u.Name, u.Age, u.Username, u.Password)
	respondWithJSON(w, http.StatusCreated, map[string]string{"status": "Пользователь добавлен"})
}

Обновление пользователя
func updateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	db.Exec("UPDATE users SET name=$1, age=$2 WHERE id=$3", u.Name, u.Age, id)
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "Данные обновлены"})
}

Удаление пользователя
func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	db.Exec("DELETE FROM users WHERE id=$1", id)
	w.WriteHeader(http.StatusNoContent)
}

2. Клиентская часть
2.1 Назначение

Клиентское приложение представляет собой консольную программу, которая:

отправляет HTTP-запросы на сервер;

выполняет авторизацию пользователей;

выполняет CRUD-операции;

поддерживает несколько пользователей с разными токенами.

2.2 Общая структура
var userTokens = make(map[string]string) // username → token
var currentUser string                   // активный пользователь


Приложение использует карту для хранения токенов всех авторизованных пользователей и позволяет переключаться между ними.

2.3 Авторизация пользователя
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

	resp, _ := http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(data))
	body, _ := io.ReadAll(resp.Body)

	var result map[string]string
	json.Unmarshal(body, &result)
	token := result["token"]

	userTokens[username] = token
	currentUser = username
}


Порядок работы:

Пользователь вводит логин и пароль.

Клиент отправляет запрос POST /login.

Сервер возвращает токен, который сохраняется в userTokens.

Этот токен будет использоваться при всех последующих запросах.

2.4 Добавление заголовка авторизации

Перед каждым запросом к серверу клиент добавляет токен текущего пользователя:

func addAuth(req *http.Request) {
	token := userTokens[currentUser]
	req.Header.Set("Authorization", "Bearer "+token)
}

2.5 CRUD операции

Примеры вызовов:

func getAllUsers() {
	req, _ := http.NewRequest("GET", baseURL+"/users", nil)
	addAuth(req)
	send(req)
}

func createUser() {
	var u User
	fmt.Print("Имя: "); fmt.Scan(&u.Name)
	fmt.Print("Возраст: "); fmt.Scan(&u.Age)
	fmt.Print("Логин: "); fmt.Scan(&u.Username)
	fmt.Print("Пароль: "); fmt.Scan(&u.Password)

	data, _ := json.Marshal(u)
	req, _ := http.NewRequest("POST", baseURL+"/users", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	addAuth(req)
	send(req)
}


Функция send(req) выполняет HTTP-запрос и выводит ответ на экран.

2.6 Поддержка нескольких пользователей

Клиент позволяет одновременно работать под разными логинами.
Токены для каждого пользователя сохраняются в map[string]string.
Пользователь может переключаться с помощью функции:

func switchUser() {
	fmt.Print("Введите имя пользователя: ")
	var name string
	fmt.Scan(&name)
	currentUser = name
}

3. Итог

Данное приложение реализует:

авторизацию пользователей с уникальными токенами;

хранение токенов в базе данных (sessions);

REST API с CRUD-операциями;

клиентскую часть с поддержкой нескольких пользователей;

валидацию и централизованную обработку ошибок.

Приложение является готовым шаблоном для построения безопасных REST API на Go с консольным клиентом.