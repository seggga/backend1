package main

import (
	"log"
	"net/http"
	"time"
)

var (
	serverAddr string = "localhost:8080"
	serveDir   string = "./files"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", uploadHandleFunc)
	mux.HandleFunc("/list", listHandleFunc)

	fs := &http.Server{
		Addr:         serverAddr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := fs.ListenAndServe()
	if err != nil {
		log.Printf("error on server, %v", err)
	}
}
