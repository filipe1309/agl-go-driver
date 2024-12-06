package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/filipe1309/agl-go-driver/application/http/users"
	"github.com/filipe1309/agl-go-driver/factories"
	"github.com/filipe1309/agl-go-driver/internal/auth"
	"github.com/filipe1309/agl-go-driver/internal/files"
	"github.com/filipe1309/agl-go-driver/internal/folders"
	"github.com/filipe1309/agl-go-driver/repositories"
	"github.com/filipe1309/agl-go-driver/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func RunHTTPServer() {
	db, bucket, queueConn := getSessions()

	r := chi.NewRouter()

	// allow cors
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	// Define endpoints
	files.SetRoutes(r, db, bucket, queueConn)
	folders.SetRoutes(r, db)

	// Users DDD
	ur := repositories.NewUserRepository(db)
	uf := factories.NewUserFactory(ur)
	users.SetRoutes(r, ur, uf)

	// Auth
	authService := services.NewAuthService(ur, uf)

	r.Post("/auth", auth.HandleHttpAuth(func(login, password string) (auth.Authenticated, error) {
		return authService.Authenticate(login, password)
	}))

	// Start server
	log.Println("Server running on port " + os.Getenv("SERVER_PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")), r)
}
