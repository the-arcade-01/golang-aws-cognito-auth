package api

import (
	"app/internal/config"
	"app/internal/db"
	"app/internal/handlers"
	"app/internal/models"
	"app/internal/services"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func newRoutes(cfg *config.Config) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(requestLogger)

	authStore, err := db.NewCognitoStore(cfg)
	if err != nil {
		slog.Error("failed to initialize auth store", "err", err)
		panic(err)
	}
	authHandlers := handlers.NewAuthHandlers(
		services.NewAuthService(
			authStore,
		),
	)
	authRouter := chi.NewRouter()

	authRouter.Post("/signup", authHandlers.SignUp)
	authRouter.Post("/login", authHandlers.Login)
	authRouter.Post("/confirm", authHandlers.ConfirmAccount)

	authRouter.Group(func(r chi.Router) {
		r.Use(jwtAuthMiddleware(authStore))
		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			models.ResponseWithJSON(w, http.StatusOK, models.NewDataResponse(http.StatusOK, "Protected route"))
		})
		r.Get("/user/info", authHandlers.GetUser)
	})

	router.Mount("/auth", authRouter)
	return router
}
