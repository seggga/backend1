package handler

import "net/http"

func (rt *Router) downloadFunc() http.Handler {
	dirToServe := http.Dir(rt.Dir)
	return http.StripPrefix("/download/", http.FileServer(dirToServe))
}
