package interceptors

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/api/grpc/meta"
	"github.com/chains-lab/distributors-svc/internal/api/grpc/problems"
	"github.com/chains-lab/gatekit/auth"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UserJwtAuth(skUser string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, problems.UnauthenticatedError(ctx, "no metadata found in incoming context")
		}

		token := md["x-user-token"]
		if len(token) == 0 {
			return nil, problems.UnauthenticatedError(ctx, fmt.Sprintf("user token not supplied"))
		}

		userData, err := auth.VerifyUserJWT(ctx, token[0], skUser)
		if err != nil {
			return nil, problems.UnauthenticatedError(ctx, "failed to verify user token")
		}

		userID, err := uuid.Parse(userData.Subject)
		if err != nil {
			return nil, problems.UnauthenticatedError(ctx, fmt.Sprintf("invalid user UserID: %v", err))
		}

		ctx = context.WithValue(ctx, meta.UserCtxKey, meta.UserData{
			ID:        userID,
			SessionID: userData.Session,
			Verified:  userData.Verified,
			Role:      userData.Role,
		})

		return handler(ctx, req)
	}
}
