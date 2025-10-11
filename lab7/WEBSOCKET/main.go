package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // разрешаем всем подключениям
}

var clients = make(map[*websocket.Conn]bool) // список подключённых клиентов
var broadcast = make(chan []byte)            // канал для рассылки сообщений

func main() {
	http.HandleFunc("/ws", handleConnections)

	go handleMessages() // отдельная горутина для рассылки сообщений

	log.Println("WebSocket сервер запущен на ws://localhost:8080/ws")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка подключения:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	log.Println("Клиент подключён")

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Клиент отключился")
			delete(clients, ws)
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("Ошибка отправки:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}