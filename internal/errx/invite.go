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

var ErrorInviteNotFound = ape.Declare("INVITE_NOT_FOUND")

func RaiseInviteNotFound(ctx context.Context, cause error, inviteID uuid.UUID) error {
	st := status.New(codes.NotFound, fmt.Sprintf("invite not found: %s", inviteID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInviteNotFound.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInviteNotFound.Raise(cause, st)
}

var ErrorInviteIsNotActive = ape.Declare("INVITE_IS_NOT_ACTIVE")

func RaiseInviteIsNotActive(ctx context.Context, cause error, inviteID uuid.UUID) error {
	st := status.New(codes.FailedPrecondition, fmt.Sprintf("invite already answered: %s", inviteID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInviteIsNotActive.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInviteIsNotActive.Raise(cause, st)
}

var ErrorInviteIsNotForInitiator = ape.Declare("INVITE_IS_NOT_FOR_INITIATOR")

func RaiseInviteIsNotForInitiator(ctx context.Context, cause error, inviteID uuid.UUID) error {
	st := status.New(codes.FailedPrecondition, fmt.Sprintf("invite is not for initiator: %s", inviteID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInviteIsNotForInitiator.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInviteIsNotForInitiator.Raise(cause, st)
}

var ErrorCantSendInviteForCurrentEmployee = ape.Declare("CANT_SEND_INVITE_FOR_CURRENT_EMPLOYEE")

func RaiseCantSendInviteForCurrentEmployee(ctx context.Context, cause error, inviteID uuid.UUID) error {
	st := status.New(codes.FailedPrecondition, fmt.Sprintf("can't send invite for current employee: %s", inviteID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorCantSendInviteForCurrentEmployee.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorCantSendInviteForCurrentEmployee.Raise(cause, st)
}

var UserHaveAlreadyInviteForInitiatorDistributor = ape.Declare("USER_HAVE_ALREADY_INVITE_FOR_INITIATOR_DISTRIBUTOR")

func RaiseUserHaveAlreadyInviteForInitiatorDistributor(ctx context.Context, cause error, distributorID uuid.UUID) error {
	st := status.New(codes.FailedPrecondition, fmt.Sprintf("user have already invite for initiator distributor: %s", distributorID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: UserHaveAlreadyInviteForInitiatorDistributor.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return UserHaveAlreadyInviteForInitiatorDistributor.Raise(cause, st)
}
