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

// ---------------------------
// Структура пользователя
// ---------------------------
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *sql.DB

// ---------------------------
// Главная функция
// ---------------------------
func main() {
	connStr := "user=postgres password=postgres dbname=usersdb sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}
	log.Println("Подключение к базе данных успешно!")

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("REST API with PostgreSQL, validation, pagination and filtering"))
	})
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

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
// Валидация пользователя
// ---------------------------
func validateUser(u User) (bool, string) {
	u.Name = strings.TrimSpace(u.Name)
	if u.Name == "" {
		return false, "Имя не может быть пустым"
	}
	if len(u.Name) > 50 {
		return false, "Имя слишком длинное (максимум 50 символов)"
	}
	for _, ch := range u.Name {
		if ch >= '0' && ch <= '9' {
			return false, "Имя не должно содержать цифр"
		}
	}
	if u.Age <= 0 || u.Age > 120 {
		return false, "Некорректный возраст"
	}
	return true, ""
}

// ---------------------------
// Обработчики маршрутов
// ---------------------------

// Получить всех пользователей с пагинацией и фильтрацией
func getUsers(w http.ResponseWriter, r *http.Request) {
	// Параметры запроса
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")
	nameFilter := r.URL.Query().Get("name")
	ageFilter := r.URL.Query().Get("age")

	// Значения по умолчанию
	page := 1
	limit := 5
	if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
		limit = l
	}
	offset := (page - 1) * limit

	// Формируем SQL-запрос с фильтрацией
	query := "SELECT id, name, age FROM users WHERE 1=1"
	var args []interface{}
	argIndex := 1

	if nameFilter != "" {
		query += " AND name ILIKE $" + strconv.Itoa(argIndex)
		args = append(args, "%"+nameFilter+"%")
		argIndex++
	}

	if ageFilter != "" {
		if age, err := strconv.Atoi(ageFilter); err == nil {
			query += " AND age = $" + strconv.Itoa(argIndex)
			args = append(args, age)
			argIndex++
		}
	}

	query += " ORDER BY id LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка выполнения запроса")
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

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"page":  page,
		"limit": limit,
		"users": users,
	})
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

	ok, msg := validateUser(u)
	if !ok {
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	err := db.QueryRow("INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id", u.Name, u.Age).Scan(&u.ID)
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

	ok, msg := validateUser(u)
	if !ok {
		respondWithError(w, http.StatusBadRequest, msg)
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
