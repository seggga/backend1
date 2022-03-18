package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/seggga/backend1/internal/api/router/defaultmux"
	"github.com/seggga/backend1/internal/api/server"
)

func main() {

	uploadDir := "upload"
	err := createFolder(uploadDir)
	if err != nil {
		log.Fatal(err)
	}

	router := defaultmux.New(uploadDir)
	srv := server.New(":8080", router)

	srv.Start()

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-osSigChan
	log.Println("received OS interrupting signal")
	srv.Stop()

}

func createFolder(s string) error {
	_, err := os.Stat(s)
	switch {
	case os.IsNotExist(err):
		log.Println("create folder", s)
		err := os.Mkdir(s, 0755)
		if err != nil {
			log.Println("cannot create folder", s, err)
			return err
		}

	default:
		return err

	}
	return nil
}
