package middleware

import (
	"net/http"
	"time"

	"github.com/tusmasoma/connectHub-backend/internal/log"
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
		log.Info(
			"Access log",
			log.Ftime("Date", time.Now()),
			log.Fstring("URL", r.URL.String()),
			log.Fstring("IP", r.RemoteAddr),
			log.Fint("StatusCode", lrw.statusCode),
		)

		// エラーログ (StatusCodeが400以上の場合)
		if lrw.statusCode >= StatusCodeBadRequest {
			log.Error(
				"Error log",
				log.Ftime("Date", time.Now()),
				log.Fstring("URL", r.URL.String()),
				log.Fint("StatusCode", lrw.statusCode),
			)
		}
	})
}
