package handler

import (
	"net/http"
)

// Router ...
type Router struct {
	*http.ServeMux
	Dir string
}

// New creates a router with specified storage and handlers
func New(dir string) *Router {
	r := &Router{
		ServeMux: http.NewServeMux(),
		Dir:      dir,
	}

	r.HandleFunc("/upload", r.uploadFunc)
	r.Handle("/download/", r.downloadFunc())
	r.HandleFunc("/list", r.listFilesFunc)
	r.HandleFunc("/filtered-list", r.filteredList)

	return r
}
