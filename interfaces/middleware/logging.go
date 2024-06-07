package middleware

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

const StatusCodeBadRequest = 400

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{ResponseWriter: w}

		next.ServeHTTP(lrw, r)
		// アクセスログ
		log.Printf("[ACCESS] Date: %s, URL: %s, IP: %s, StatusCode: %d",
			time.Now().Format("2006-01-02 15:04:05"), r.URL, r.RemoteAddr, lrw.statusCode)

		// エラーログ (StatusCodeが400以上の場合)
		if lrw.statusCode >= StatusCodeBadRequest {
			log.Printf("[ERROR] Date: %s, URL: %s, StatusCode: %d",
				time.Now().Format("2006-01-02 15:04:05"), r.URL, lrw.statusCode)
		}
	})
}
