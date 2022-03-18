package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (h *Hands) UploadFunc(w http.ResponseWriter, r *http.Request) {

	log.Println("/upload handler has been called")

	if r.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

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

	filePath := h.dir + "/" + header.Filename
	log.Println("writing file ", filePath)
	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	fileLink := "/download/" + header.Filename
	fmt.Fprintln(w, fileLink)
}
