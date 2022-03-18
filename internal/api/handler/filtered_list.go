package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func (h *Hands) FilteredList(w http.ResponseWriter, r *http.Request) {

	log.Println("/filtered-list handler has been called")

	if r.Method != http.MethodGet {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	filter := r.FormValue("extension")
	log.Println("filter by extension", filter)
	if filter == "" {
		http.Error(w, "Missed extension", http.StatusBadRequest)
		return
	}

	files, err := ioutil.ReadDir(h.dir)
	if err != nil {
		log.Fatal(err)
	}

	encoder := json.NewEncoder(w)
	w.WriteHeader(http.StatusOK)
	for _, f := range files {
		name := f.Name()
		fileExtension := filepath.Ext(name)
		if fileExtension != filter {
			continue
		}

		_ = encoder.Encode(
			// DTO
			File{
				Name:        name,
				Extension:   fileExtension,
				SizeInBytes: int(f.Size()),
				FileLink:    "/download/" + name,
			},
		)
	}
}
