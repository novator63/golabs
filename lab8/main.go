package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *sql.DB

func main() {
	// Подключение к базе данных
	connStr := "user=postgres password=postgres dbname=usersdb sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверяем соединение
	if err = db.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}
	log.Println("Подключение к базе данных успешно!")

	// Настройка маршрутов
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("REST API with PostgreSQL is running"))
	})
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	// Запуск сервера
	log.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", r)
}

// ---------------------------
// Универсальные функции ответов
// ---------------------------
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// ---------------------------
// Обработчики маршрутов
// ---------------------------

// Получить всех пользователей
func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка базы данных")
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Ошибка чтения данных")
			return
		}
		users = append(users, u)
	}
	respondWithJSON(w, http.StatusOK, users)
}

// Получить пользователя по ID
func getUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var u User
	err := db.QueryRow("SELECT id, name, age FROM users WHERE id=$1", id).Scan(&u.ID, &u.Name, &u.Age)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Пользователь не найден")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка базы данных")
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

// Добавить нового пользователя
func createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}

	// Валидация
	u.Name = strings.TrimSpace(u.Name)
	if u.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Имя не может быть пустым")
		return
	}
	if u.Age <= 0 {
		respondWithError(w, http.StatusBadRequest, "Возраст должен быть положительным числом")
		return
	}

	err := db.QueryRow(
		"INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id",
		u.Name, u.Age,
	).Scan(&u.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка вставки в базу")
		return
	}

	respondWithJSON(w, http.StatusCreated, u)
}

// Обновить пользователя
func updateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}

	// Валидация
	u.Name = strings.TrimSpace(u.Name)
	if u.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Имя не может быть пустым")
		return
	}
	if u.Age <= 0 {
		respondWithError(w, http.StatusBadRequest, "Возраст должен быть положительным числом")
		return
	}

	res, err := db.Exec("UPDATE users SET name=$1, age=$2 WHERE id=$3", u.Name, u.Age, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка базы данных")
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		respondWithError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}
	u.ID = id
	respondWithJSON(w, http.StatusOK, u)
}

// Удалить пользователя
func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	res, err := db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка базы данных")
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		respondWithError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
