package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {

	// read nickname from stdin
	fmt.Print("Enter your nickname:")
	reader := bufio.NewReader(os.Stdin)
	nickname, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println("cannot read data, program exit")
		return
	}

	// establish connection
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// send nickname to server
	_, err = conn.Write(nickname)
	if err != nil {
		fmt.Printf("could not send nickname to server, %v", err)
		return
	}

	// listen for the data from server
	go func() {
		io.Copy(os.Stdout, conn)
	}()

	// send data to the server
	io.Copy(conn, os.Stdin) // until you send ^Z
	fmt.Printf("%s: exit", conn.LocalAddr())
}
