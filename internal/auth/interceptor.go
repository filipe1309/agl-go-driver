package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ValidateTokenInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if len(md["authorization"]) > 0 {
			token := md["authorization"][0]
			// Authorization: Bearer <token>
			token = strings.TrimPrefix(token, "Bearer ")

			claims, err, _ := validate(token)
			if err != nil {
				return nil, err
			}

			ctx := context.WithValue(ctx, "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "user_name", claims.UserName)

			return handler(ctx, req)
		}
	}

	return nil, status.Errorf(codes.Unauthenticated, "missing token")
}
