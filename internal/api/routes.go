package api

import (
	"app/internal/config"
	"app/internal/db"
	"app/internal/handlers"
	"app/internal/services"
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

	authHandlers := handlers.NewAuthHandlers(
		services.NewAuthService(
			db.NewCognitoStore(cfg),
		),
	)
	authRouter := chi.NewRouter()

	authRouter.Post("/signup", authHandlers.SignUp)
	authRouter.Post("/login", authHandlers.Login)
	authRouter.Post("/confirm", authHandlers.ConfirmAccount)

	router.Mount("/auth", authRouter)
	return router
}
