package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GrpcAuthInterceptor returns a unary server interceptor that validates JWT tokens.
func GrpcAuthInterceptor(secret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "Missing metadata")
		}

		tokenStr := extractToken(md)
		if tokenStr == "" {
			return nil, status.Errorf(codes.Unauthenticated, "Authorization token required")
		}

		token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid or expired token")
		}

		// Proceed to actual RPC
		return handler(ctx, req)
	}
}

func extractToken(md metadata.MD) string {
	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return ""
	}
	parts := strings.SplitN(authHeaders[0], " ", 2)
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return ""
}
