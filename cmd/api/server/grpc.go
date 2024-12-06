package server

import (
	"log"
	"net"

	"github.com/filipe1309/agl-go-driver/application/grpc/users"
	"github.com/filipe1309/agl-go-driver/factories"
	"github.com/filipe1309/agl-go-driver/internal/auth"
	authpb "github.com/filipe1309/agl-go-driver/proto/v1/auth"
	userspb "github.com/filipe1309/agl-go-driver/proto/v1/users"
	"github.com/filipe1309/agl-go-driver/repositories"
	"github.com/filipe1309/agl-go-driver/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func RunGRPCServer() {
	// Start grpc server
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	db, _, _ := getSessions()

	// Users DDD
	ur := repositories.NewUserRepository(db)
	uf := factories.NewUserFactory(ur)
	userService := users.NewUserService(ur, uf)

	// Auth
	authService := services.NewAuthService(ur, uf)

	// svc and method bypass token validation
	auth.AddByPassValidateToken("auth.AuthService", "Login")
	auth.AddByPassValidateToken("users.UserService", "Create")

	grpcAuthService := auth.HandleGrpcAuth(func(login, password string) (auth.Authenticated, error) {
		return authService.Authenticate(login, password)
	})

	creds, err := credentials.NewServerTLSFromFile("certs/server/server.crt", "certs/server/server.key")
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}

	// Define grpc server
	s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(auth.ValidateTokenInterceptor))
	authpb.RegisterAuthServiceServer(s, grpcAuthService)
	userspb.RegisterUserServiceServer(s, userService)

	log.Println("Server running on port 50051")
	if err := s.Serve(l); err != nil {
		panic(err)
	}
}
