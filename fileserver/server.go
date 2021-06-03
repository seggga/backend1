package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/seggga/backend1/fileserver/handlers"
)

func main() {
	handler := &handlers.Handler{}
	uploadHandler := &handlers.UploadHandler{
		UploadDir: "upload",
	}

	http.Handle("/", handler)
	http.Handle("/upload", uploadHandler)

	srv := &http.Server{
		Addr:         ":80",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	dirToServe := http.Dir(uploadHandler.UploadDir)
	fs := &http.Server{
		Addr:         ":8080",
		Handler:      http.FileServer(dirToServe),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Println("error starting server")
			return
		}
	}()

	go func() {
		err := fs.ListenAndServe()
		if err != nil {
			fmt.Println("error starting server")
			return
		}
	}()

}
