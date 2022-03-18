package defaultmux

import (
	"net/http"

	"github.com/seggga/backend1/internal/api/handler"
)

// Router ...
type Router struct {
	*http.ServeMux
	Dir string
	*handler.Hands
}

type hands interface {
}

// New creates a router with specified storage and handlers
func New(dir string) *Router {
	r := &Router{
		ServeMux: http.NewServeMux(),
		Dir:      dir,
		Hands:    handler.New(dir),
	}

	r.HandleFunc("/upload", r.UploadFunc)
	r.Handle("/download/", r.DownloadFunc())
	r.HandleFunc("/list", r.ListFilesFunc)
	r.HandleFunc("/filtered-list", r.FilteredList)

	return r
}
