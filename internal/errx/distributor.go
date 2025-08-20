package errx

import (
	"context"

	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorDistributorNotFound = ape.Declare("DISTRIBUTOR_NOT_FOUND")

func RaiseDistributorNotFound(ctx context.Context, cause error, distributorID uuid.UUID) error {
	st := status.New(codes.NotFound, "distributor not found")
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorDistributorNotFound.Error(),
			Domain: "city-petitions-svc",
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorDistributorNotFound.Raise(cause, st)
}
