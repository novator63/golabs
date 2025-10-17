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

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
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

	// Публичный маршрут — логин
	r.HandleFunc("/login", loginHandler).Methods("POST")

	// Защищённые маршруты
	api := r.PathPrefix("/").Subrouter()
	api.Use(authMiddleware)
	api.HandleFunc("/users", getUsers).Methods("GET")
	api.HandleFunc("/users/{id}", getUser).Methods("GET")
	api.HandleFunc("/users", createUser).Methods("POST")
	api.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	log.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", r)
}

// ---------------------------
// Авторизация и middleware
// ---------------------------
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}

	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username=$1 AND password=$2",
		creds.Username, creds.Password).Scan(&userID)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusUnauthorized, "Неверные учетные данные")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка базы данных")
		return
	}

	// Генерация уникального токена
	token := uuid.New().String()

	_, err = db.Exec("INSERT INTO sessions (user_id, token, created_at) VALUES ($1, $2, $3)",
		userID, token, time.Now())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка сохранения токена")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

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
		} else if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Ошибка проверки токена")
			return
		}

		next.ServeHTTP(w, r)
	})
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
// CRUD
// ---------------------------
func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, age, username FROM users ORDER BY id")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка чтения из базы")
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Age, &u.Username)
		users = append(users, u)
	}
	respondWithJSON(w, http.StatusOK, users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var u User
	err := db.QueryRow("SELECT id, name, age, username FROM users WHERE id=$1", id).Scan(&u.ID, &u.Name, &u.Age, &u.Username)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	json.NewDecoder(r.Body).Decode(&u)

	_, err := db.Exec("INSERT INTO users (name, age, username, password) VALUES ($1, $2, $3, $4)",
		u.Name, u.Age, u.Username, u.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка вставки пользователя")
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"status": "Пользователь добавлен"})
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var u User
	json.NewDecoder(r.Body).Decode(&u)

	_, err := db.Exec("UPDATE users SET name=$1, age=$2 WHERE id=$3", u.Name, u.Age, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка обновления")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "Данные обновлены"})
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка удаления")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
