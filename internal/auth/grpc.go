package auth

import (
	"context"

	authpb "github.com/filipe1309/agl-go-driver/proto/v1/auth"
)

type ServiceGRPC struct {
	authpb.UnimplementedAuthServiceServer
	authHandler handler
}

func (svc *ServiceGRPC) Login(ctx context.Context, req *authpb.CredentialsRequest) (*authpb.TokenResponse, error) {
	credentials := Credentials{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}
	token, err, _ := svc.authHandler.auth(credentials)
	if err != nil {
		return &authpb.TokenResponse{Error: err.Error()}, err
	}

	return &authpb.TokenResponse{Token: token}, nil
}

func HandleGrpcAuth(fn authenticateFunc) *ServiceGRPC {
	svc := &ServiceGRPC{
		authHandler: handler{fn},
	}

	return svc
}
