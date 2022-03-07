package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

func main() {
	handler := &Handler{
		HostAddr:  ":8080",
		UploadDir: "upload",
	}

	http.HandleFunc("/upload", handler.uploadFunc)
	http.Handle("/download/", handler.downloadFunc())
	http.HandleFunc("/list", handler.listFilesFunc)

	fs := &http.Server{
		Addr:         handler.HostAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("server has been started on address ", handler.HostAddr)
	fs.ListenAndServe()
}

// Handler ...
type Handler struct {
	HostAddr  string
	UploadDir string
}

func (h *Handler) uploadFunc(w http.ResponseWriter, r *http.Request) {

	log.Println("/upload handler has been called")

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Println("got file ", header.Filename)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		log.Println("Unable to read file")
		return
	}

	filePath := h.UploadDir + "/" + header.Filename
	log.Println("writing file ", filePath)
	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	fileLink := h.HostAddr + "/download/" + header.Filename
	fmt.Fprintln(w, fileLink)
}

func (h *Handler) downloadFunc() http.Handler {
	dirToServe := http.Dir(h.UploadDir)
	return http.StripPrefix("/download/", http.FileServer(dirToServe))
}

// File struct represents a file returned as response
type File struct {
	Name        string
	Extention   string
	SizeInBytes int
	FileLink    string
}

func (h *Handler) listFilesFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("/list handler has been called")

	files, err := ioutil.ReadDir(h.UploadDir)
	if err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	for _, f := range files {
		name := f.Name()
		fileExtension := filepath.Ext(name)

		_ = encoder.Encode(
			// DTO
			File{
				Name:        name,
				Extention:   fileExtension,
				SizeInBytes: int(f.Size()),
				FileLink:    h.HostAddr + "/download/" + name,
			},
		)
	}
}
