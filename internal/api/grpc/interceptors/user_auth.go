package interceptors

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/api/grpc/meta"
	"github.com/chains-lab/distributors-svc/internal/errx"
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
			return nil, errx.RaiseUnauthenticated(ctx, fmt.Errorf("no metadata found in incoming context"))
		}

		token := md["x-user-token"]
		if len(token) == 0 {
			return nil, errx.RaiseUnauthenticated(ctx, fmt.Errorf("user token not supplied"))
		}

		userData, err := auth.VerifyUserJWT(ctx, token[0], skUser)
		if err != nil {
			return nil, errx.RaiseUnauthenticated(ctx, fmt.Errorf("failed to verify user token"))
		}

		userID, err := uuid.Parse(userData.Subject)
		if err != nil {
			return nil, errx.RaiseUnauthenticated(ctx, fmt.Errorf("invalid user UserID: %v", err))
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
