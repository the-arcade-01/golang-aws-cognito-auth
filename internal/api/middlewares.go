package api

import (
	"log/slog"
	"net/http"
	"time"
)

type statusRecorder struct {
	status int
	http.ResponseWriter
}

func (s *statusRecorder) WriteHeader(status int) {
	s.status = status
	s.ResponseWriter.WriteHeader(status)
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w}
		next.ServeHTTP(rec, r)

		micro := time.Since(start).Microseconds()
		slog.Info("rtt",
			"method", r.Method,
			"url", r.URL.String(),
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"micro", micro,
			"status", rec.status,
		)
	})
}
