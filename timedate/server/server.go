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

// broadcaster controls clients's connections and sends them messages
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

// msgGenerator produces a message once a second
func msgGenerator() {
	// channel for user's messages to compete with 1-second ticker
	userChan := make(chan string)
	// read messages from server's console
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter text: ")
			text, _ := reader.ReadString('\n')
			text = text[:len(text)-1] // crop the last symbol '\n'
			userChan <- text
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	var message string
	for {
		// send a message once a second
		<-ticker.C
		select {
		case userMsg := <-userChan: // user's message has come
			message = time.Now().Format("15:04:05") + " " + userMsg
		default: // no messages has come
			message = time.Now().Format("15:04:05")
		}
		messages <- message
	}
}

// handleCon works with patrticular client's connection
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
