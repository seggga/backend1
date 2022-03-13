package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/seggga/backend1/internal/api/handler"
	"github.com/seggga/backend1/internal/api/server"
)

func main() {
	handler := handler.New("upload")
	srv := server.New(":8080", handler)

	srv.Start()

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-osSigChan
	log.Println("received OS interrupting signal")
	srv.Stop()

}
