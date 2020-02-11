package main

import (
	"log"
	"net/http"
)

func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		rw := responseWriter{
			ResponseWriter: res,
		}
		h.ServeHTTP(&rw, req)
		log.Println("[http]",
			req.Method, req.URL.String(),
			rw.status, http.StatusText(rw.status),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	return rw.ResponseWriter.Write(p)
}
