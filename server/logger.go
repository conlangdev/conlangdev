package server

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	if r.responseData.status == 0 {
		r.responseData.status = 200
	}
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Inject logging response writer into server stack
		rd := &responseData{0, 0}
		lrw := loggingResponseWriter{w, rd}
		next.ServeHTTP(&lrw, r)
		// Log what we found out!
		log.WithFields(log.Fields{
			"method":   r.Method,
			"uri":      r.RequestURI,
			"status":   rd.status,
			"size":     rd.size,
			"duration": time.Since(start),
		}).Info("Request complete.")
	})
}
