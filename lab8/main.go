package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Структура пользователя
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users []User
var nextID = 1

func main() {
	r := mux.NewRouter()

	// Маршруты
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	http.ListenAndServe(":8080", r)
}

// Получить список всех пользователей
func getUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

// Получить пользователя по ID
func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for _, user := range users {
		if user.ID == id {
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

// Добавить нового пользователя
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	user.ID = nextID
	nextID++
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
}

// Обновить пользователя по ID
func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i, user := range users {
		if user.ID == id {
			json.NewDecoder(r.Body).Decode(&users[i])
			users[i].ID = id
			json.NewEncoder(w).Encode(users[i])
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

// Удалить пользователя по ID
func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}
