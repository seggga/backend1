package main

import (
	"net/http"
	"time"
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
	srv.ListenAndServe()
}
