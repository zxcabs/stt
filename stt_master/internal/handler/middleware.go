package handler

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (r *responseWriter) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wr := NewResponseWriter(w)
		next.ServeHTTP(wr, r)

		log.Printf("(%d)[%s] %s %s [%v]", wr.statusCode, r.RemoteAddr, r.Method, r.URL.Path, time.Since(start))
	})
}
