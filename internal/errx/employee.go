package errx

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

var ErrorEmployeeNotFound = ape.Declare("EMPLOYEE_NOT_FOUND")

func RaiseEmployeeNotFound(ctx context.Context, cause error, userID uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("employee not found: %s", userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorEmployeeNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorEmployeeNotFound.Raise(cause, st)
}

var ErrorEmployeeAlreadyExists = ape.Declare("EMPLOYEE_ALREADY_EXISTS")

func RaiseEmployeeAlreadyExists(ctx context.Context, cause error, userID uuid.UUID) error {
	st := status.New(codes.AlreadyExists, fmt.Sprintf("employee already exists: %s", userID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorEmployeeAlreadyExists.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorEmployeeAlreadyExists.Raise(cause, st)
}

var ErrorInitiatorNotEmployee = ape.Declare("INITIATOR_NOT_EMPLOYEE")

func RaiseInitiatorNotEmployee(ctx context.Context, cause error, userID, distributorID uuid.UUID) error {
	st := status.New(
		codes.FailedPrecondition,
		fmt.Sprintf("user %s is not employee in distributor %s", userID, distributorID),
	)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInitiatorNotEmployee.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInitiatorNotEmployee.Raise(cause, st)
}

var ErrorInitiatorEmployeeHaveNotEnoughPermissions = ape.Declare("INITIATOR_EMPLOYEE_HAVE_NOT_ENOUGH_PERMISSIONS")

func RaiseInitiatorEmployeeHaveNotEnoughPermissions(ctx context.Context, cause error, userID, distributorID uuid.UUID) error {
	st := status.New(
		codes.PermissionDenied,
		fmt.Sprintf("employee %s have not eniugh right in distributor %s", userID, distributorID),
	)
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInitiatorEmployeeHaveNotEnoughPermissions.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInitiatorEmployeeHaveNotEnoughPermissions.Raise(cause, st)
}
