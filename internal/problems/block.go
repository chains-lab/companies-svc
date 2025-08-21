package problems

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/config/constant"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorBlockNotFound = ape.Declare("LOCKOUT_NOT_FOUND")

func RaiseBlockNotFound(ctx context.Context, cause error, ID uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("block %s not found", ID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorBlockNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorBlockNotFound.Raise(cause, st)
}

var DistributorHaveNotActiveBlock = ape.Declare("DISTRIBUTOR_HAVE_NOT_ACTIVE_LOCKOUT")

func RaiseDistributorHaveNotActiveBlock(ctx context.Context, cause error, distributorID uuid.UUID) error {
	st := status.New(codes.FailedPrecondition, fmt.Sprintf("distributor %s have not active block", distributorID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: DistributorHaveNotActiveBlock.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return DistributorHaveNotActiveBlock.Raise(cause, st)
}

var DistributorHaveAlreadyActiveBlock = ape.Declare("DISTRIBUTOR_HAVE_ALREADY_ACTIVE_LOCKOUT")

func RaiseDistributorHaveAlreadyActiveBlock(ctx context.Context, cause error, distributorID uuid.UUID) error {
	st := status.New(codes.FailedPrecondition, fmt.Sprintf("distributor %s have already active block", distributorID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: DistributorHaveAlreadyActiveBlock.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return DistributorHaveAlreadyActiveBlock.Raise(cause, st)
}
