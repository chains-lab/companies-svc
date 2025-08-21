package interceptors

import (
	"context"

	"github.com/chains-lab/distributors-svc/internal/api/grpc/meta"
	"github.com/chains-lab/gatekit/roles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func AuthAdmin() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		switch info.FullMethod {
		case "TODO method name":
			user, err := meta.User(ctx)
			if err != nil {
				return nil, err
			}

			if user.Role != roles.Admin && user.Role != roles.SuperUser {
				return nil, status.Error(codes.PermissionDenied, "user does not have admin permissions")
			}
		}

		return handler(ctx, req)
	}
}
