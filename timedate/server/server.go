package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	go msgGenerator()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func msgGenerator() {
	// channel for user's messages to compete with 1-second ticker
	userChan := make(chan string)

	// read messages from server's console
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			userChan <- scanner.Text()
		}

	}()

	ticker := time.NewTicker(1 * time.Second)

	for {

		<-ticker.C
		select {
		case userMsg := <-userChan:
			messages <- time.Now().Format("15:04:05\n\r") + userMsg
		default:
			messages <- time.Now().Format("15:04:05\n\r")
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)

	entering <- ch

	for msg := range ch {
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			break
		}
	}

	leaving <- ch
	conn.Close()
}
