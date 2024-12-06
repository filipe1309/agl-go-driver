package requests

import (
	"log"

	"google.golang.org/grpc"
)

func GetGRPCConn() *grpc.ClientConn {
	conn, err := grpc.Dial(":50051")
	if err != nil {
		log.Fatalf("Failed to dial gRPC server: %v", err)
	}

	return conn
}
