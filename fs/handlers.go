package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

// uploader is http.HandleFunc that saves files in the serveDir
func uploader(w http.ResponseWriter, r *http.Request) {
	// get file from the form named "file" from the request
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Print("cannot read file")
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// read the file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}

	// write file on server
	filePath := serveDir + "/" + header.Filename
	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	// send back a link of saved file
	fileLink := serverAddr + "/" + header.Filename
	fmt.Fprintln(w, fileLink)

}

// lister is http.HandleFunc that lists data about files stored in serveDir
func lister(w http.ResponseWriter, r *http.Request) {
	// check for GET method
	if r.Method != http.MethodGet {
		http.Error(w, "GET-method expected", http.StatusBadRequest)
		return
	}

	// obtain list of files
	files, err := ioutil.ReadDir(serveDir)
	if err != nil {
		log.Printf("cannot read file from directory %s", serveDir)
		http.Error(w, "Unable to read directory", http.StatusBadRequest)
		return
	}

	var filesData string

	for _, someFile := range files {
		if !someFile.IsDir() {
			filesData += someFile.Name()
			filesData += "\t"
			filesData += filepath.Ext(someFile.Name())
			filesData += "\t"
			filesData += strconv.Itoa(int(someFile.Size()))
			filesData += "\n"
		}
	}
	fmt.Fprintln(w, filesData)
}
