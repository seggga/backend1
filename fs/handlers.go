package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

type FileData struct {
	Name string
	Ext  string
	Size int64
}

// uploader is http.HandleFunc that saves files in the serveDir
func uploadHandleFunc(w http.ResponseWriter, r *http.Request) {
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
	fileLink := "http://" + serverAddr + "/" + header.Filename
	fmt.Fprintln(w, fileLink)

}

// lister is http.HandleFunc that lists data about files stored in serveDir
func listHandleFunc(w http.ResponseWriter, r *http.Request) {
	// check if method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "GET-method expected", http.StatusBadRequest)
		return
	}

	// obtain list of filesystem objects on server
	files, err := ioutil.ReadDir(serveDir)
	if err != nil {
		log.Printf("cannot read file from directory %s", serveDir)
		http.Error(w, "Unable to read directory", http.StatusBadRequest)
		return
	}

	// get extension filter from query
	ext := r.FormValue("ext")
	// find files and make a slice of file-attributes
	var filesData []FileData
	for _, someFile := range files {
		if !someFile.IsDir() {
			fileAttr := FileData{
				Name: someFile.Name(),
				Ext:  filepath.Ext(someFile.Name()),
				Size: someFile.Size(),
			}
			// filter (if set) output
			if ext == "" || fileAttr.Ext == ext {
				filesData = append(filesData, fileAttr)
			}
		}
	}

	// create a response JSON
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(filesData)
	if err != nil {
		log.Printf("cannot convert data into json %v", err)
		http.Error(w, "Unable to read directory", http.StatusBadRequest)
		return
	}
}
