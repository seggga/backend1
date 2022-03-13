package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func (rt *Router) listFilesFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("/list handler has been called")

	if r.Method != http.MethodGet {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	files, err := ioutil.ReadDir(rt.Dir)
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
				FileLink:    "/download/" + name,
			},
		)
	}
}
