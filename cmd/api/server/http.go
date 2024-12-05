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
	"github.com/go-chi/chi/v5"
)

func RunHTTPServer() {
	db, bucket, queueConn := getSessions()

	// Define endpoints
	r := chi.NewRouter()
	r.Post("/auth", auth.NewHandlerAuth(func(login, password string) (auth.Authenticated, error) {
		return users.Authenticate(login, password)
	}))

	files.SetRoutes(r, db, bucket, queueConn)
	folders.SetRoutes(r, db)

	// Users DDD
	ur := repositories.NewUserRepository(db)
	uf := factories.NewUserFactory(ur)
	users.SetRoutes(r, ur, uf)

	// Start server
	log.Println("Server running on port " + os.Getenv("SERVER_PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")), r)
}
