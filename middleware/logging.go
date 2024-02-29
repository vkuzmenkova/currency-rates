package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggingRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &LoggingResponseWriter{w, 0}

		next.ServeHTTP(lrw, r)

		log.Info().Msgf("[%s] %d %s %s - %v\n", r.Method, lrw.StatusCode, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}
