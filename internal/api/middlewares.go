package api

import (
	"app/internal/db"
	"app/internal/models"
	"context"
	"log/slog"
	"net/http"
	"strings"
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

func jwtAuthMiddleware(authStore db.AuthStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				models.ResponseWithJSON(w, http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, "Authorization header required"))
				return
			}

			parts := strings.Split(header, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				models.ResponseWithJSON(w, http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, "Invalid authorization header format"))
				return
			}

			token, err := authStore.ValidateToken(parts[1])
			if err != nil {
				models.ResponseWithJSON(w, http.StatusUnauthorized, models.NewErrorResponse(http.StatusUnauthorized, "Invalid authorization token "+err.Error()))
				return
			}

			userInfo, err := authStore.GetClaims(token)
			if err != nil {
				models.ResponseWithJSON(w, http.StatusInternalServerError, models.NewErrorResponse(http.StatusInternalServerError, "Failed to extract user info"))
				return
			}

			ctx := context.WithValue(r.Context(), models.RequestContextKey, &models.RequestContext{
				UserInfo: userInfo,
				Token:    parts[1],
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
