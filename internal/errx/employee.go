package errx

import (
	"context"

	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorEmployeeNotFound = ape.Declare("EMPLOYEE_NOT_FOUND")

func RaiseEmployeeNotFound(ctx context.Context, cause error, userID, distributorID uuid.UUID) error {
	st := status.New(codes.NotFound, "employee not found")
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorEmployeeNotFound.Error(),
			Domain: "city-petitions-svc",
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorEmployeeNotFound.Raise(cause, st)
}

var ErrorEmployeeAlreadyExists = ape.Declare("EMPLOYEE_ALREADY_EXISTS")

func RaiseEmployeeAlreadyExists(ctx context.Context, cause error, userID uuid.UUID) error {
	st := status.New(codes.AlreadyExists, "employee already exists")
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorEmployeeAlreadyExists.Error(),
			Domain: "city-petitions-svc",
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorEmployeeAlreadyExists.Raise(cause, st)
}

var ErrorInitiatorNotEmployee = ape.Declare("INITIATOR_NOT_EMPLOYEE")

func RaiseInitiatorNotEmployee(ctx context.Context, cause error, userID, distributorID uuid.UUID) error {
	st := status.New(codes.FailedPrecondition, "initiator is not an employee")
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInitiatorNotEmployee.Error(),
			Domain: "city-petitions-svc",
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInitiatorNotEmployee.Raise(cause, st)
}

var ErrorInitiatorEmployeeHaveNotEnoughPermissions = ape.Declare("INITIATOR_EMPLOYEE_HAVE_NOT_ENOUGH_PERMISSIONS")

func RaiseInitiatorEmployeeHaveNotEnoughPermissions(ctx context.Context, cause error, userID, distributorID uuid.UUID) error {
	st := status.New(codes.PermissionDenied, "initiator employee have not enough permissions")
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInitiatorEmployeeHaveNotEnoughPermissions.Error(),
			Domain: "city-petitions-svc",
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInitiatorEmployeeHaveNotEnoughPermissions.Raise(cause, st)
}
