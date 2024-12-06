package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var bypass = map[string]string{}

func AddByPassValidateToken(srv, method string) {
	bypass[srv] = method
}

func ValidateTokenInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// bypass token validation
	parts := strings.Split(info.FullMethod, "/")
	if len(parts) == 3 {
		// package.service/method
		svcName := parts[1]
		methodName := parts[2]

		if v, ok := bypass[svcName]; ok {
			if v == methodName {
				return handler(ctx, req)
			}
		}
	}

	// validate token
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
