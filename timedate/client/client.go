package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		_, err = io.Copy(os.Stdout, conn)
		if err != nil {
			break
		}
	}
}
