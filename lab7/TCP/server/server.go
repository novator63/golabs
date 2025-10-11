package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
)

func main() {
	addr := ":9000"
	ln, err := net.Listen("tcp", addr) //запуск и прослушивание порта
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	log.Printf("listening on %s", addr)

	stop := make(chan os.Signal, 1) //канал для получения сигнала завершения
	signal.Notify(stop, os.Interrupt)

	go func() {
		<-stop
		log.Println("получен сигнал завершения — закрываю listener...")
		_ = ln.Close() // размораживаем Accept
	}()

	var wg sync.WaitGroup

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) { //при закрытии listener, проваливается в ошибку и выходит из
				log.Println("listener закрыт, выхожу из accept loop")
				break
			}
			log.Printf("accept error: %v", err)
			continue
		}

		wg.Add(1)
		go func(c net.Conn) { //создаём по горутине для каждого клиента
			defer wg.Done()
			handleConn(c)
		}(conn)
	}

	log.Println("ожидаю завершения активных клиентов...")
	wg.Wait()
	log.Println("сервер остановлен корректно.")
}

func handleConn(conn net.Conn) {
	log.Printf("client connected: %s", conn.RemoteAddr())
	defer func() {
		_ = conn.Close()
		log.Printf("client disconnected: %s", conn.RemoteAddr())
	}()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		clientMessage := scanner.Text()
		fmt.Printf("Received from client: %s\n", clientMessage)

		if _, err := fmt.Fprintln(conn, "Message received."); err != nil {
			log.Printf("write error to %s: %v", conn.RemoteAddr(), err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("read error from %s: %v", conn.RemoteAddr(), err)
	}
}
