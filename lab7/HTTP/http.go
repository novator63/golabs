package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("%s %s %v", r.Method, r.URL.Path, duration)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Write([]byte("Hello, world!"))
	})

	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close() //r.Body - тело HTTP-запроса, поток, который нужно закрывать, чтобы освободить ресурсы
		var payload any

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid  JSON", http.StatusBadRequest)
			return
		}

		pretty, _ := json.MarshalIndent(payload, "", " ")
		log.Printf("Получены данные:\n%s\n", pretty)

		w.WriteHeader(http.StatusOK)
	})

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", logMiddleware(mux)))
}
