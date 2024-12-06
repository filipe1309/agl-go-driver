package requests

import (
	"context"
	"log"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	grpccreds "google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var grpcAddr = ":50051"

func GetGRPCConn() *grpc.ClientConn {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(getCreds()))
	if err != nil {
		log.Fatalf("Failed to dial gRPC server: %v", err)
	}

	return conn
}

func GetGRPCWithTokenConn() *grpc.ClientConn {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(getCreds()), grpc.WithUnaryInterceptor(tokenInterceptor))
	if err != nil {
		log.Fatalf("Failed to dial gRPC server: %v", err)
	}

	return conn
}

func getCreds() grpccreds.TransportCredentials {
	creds, err := grpccreds.NewClientTLSFromFile("certs/ca.crt", "example.com")
	if err != nil {
		log.Fatalf("Failed to create TLS credentials: %v", err)
	}

	return creds
}

func tokenInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	token, err := readCacheToken()
	if err != nil {
		log.Println("No token found, please login")
		return err
	}

	opts = append(opts, grpc.PerRPCCredentials(oauth.TokenSource{
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
	}))

	return invoker(ctx, method, req, reply, cc, opts...)
}
