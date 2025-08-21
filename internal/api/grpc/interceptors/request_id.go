package interceptors

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/api/grpc/meta"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RequestID() grpc.UnaryServerInterceptor {
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

		requestIDArr := md["x-request-id"]
		if len(requestIDArr) == 0 {
			return nil, errx.RaiseUnauthenticated(ctx, fmt.Errorf("request UserID not supplied"))
		}

		requestID, err := uuid.Parse(requestIDArr[0])
		if err != nil {
			return nil, errx.RaiseUnauthenticated(ctx, fmt.Errorf("invalid request UserID: %v", err))
		}

		ctx = context.WithValue(ctx, meta.RequestIDCtxKey, requestID)

		return handler(ctx, req)
	}
}
