package handler

import "net/http"

func (h *Hands) DownloadFunc() http.Handler {
	dirToServe := http.Dir(h.dir)
	return http.StripPrefix("/download/", http.FileServer(dirToServe))
}
