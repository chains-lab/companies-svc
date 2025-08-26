package problems

import (
	"context"
	"fmt"

	"github.com/chains-lab/distributors-svc/internal/config/constant"
	"github.com/chains-lab/distributors-svc/pkg/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorDistributorNotFound = ape.Declare("DISTRIBUTOR_NOT_FOUND")

func RaiseDistributorNotFound(ctx context.Context, cause error, distributorID uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("distributor %s not found", distributorID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorDistributorNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorDistributorNotFound.Raise(cause, st)
}

var ErrorDistributorStatusBlocked = ape.Declare("DISTRIBUTOR_STATUS_BLOCKED")

func RaiseDistributorStatusBlocked(ctx context.Context, cause error, distributorID uuid.UUID) error {
	st := status.New(codes.FailedPrecondition, fmt.Sprintf("distributor %s is blocked", distributorID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorDistributorStatusBlocked.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorDistributorStatusBlocked.Raise(cause, st)
}

var ErrorCurrentEmployeeCanNotCreateDistributor = ape.Declare("CURRENT_EMPLOYEE_CAN_NOT_CREATE_DISTRIBUTOR")

func RaiseCurrentEmployeeCanNotCreateDistributor(ctx context.Context, cause error, userID uuid.UUID) error {
	st := status.New(codes.FailedPrecondition, fmt.Sprintf("current employee %s can not create distributor", userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCurrentEmployeeCanNotCreateDistributor.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorCurrentEmployeeCanNotCreateDistributor.Raise(cause, st)
}
