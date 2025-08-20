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

var ErrorInviteAlreadyExists = ape.Declare("INVITE_ALREADY_EXISTS")

func RaiseInviteAlreadyExists(ctx context.Context, cause error, inviteID uuid.UUID) error {
	st := status.New(codes.AlreadyExists, fmt.Sprintf("invite already exists: %s", inviteID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInviteAlreadyExists.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInviteAlreadyExists.Raise(cause, st)
}

var ErrorInviteAlreadyAnswered = ape.Declare("INVITE_ALREADY_ANSWERED")

func RaiseInviteAlreadyAnswered(ctx context.Context, cause error, inviteID uuid.UUID) error {
	st := status.New(codes.FailedPrecondition, fmt.Sprintf("invite already answered: %s", inviteID))
	st, _ = st.WithDetails(
		&errdetails.ErrorInfo{
			Reason: ErrorInviteAlreadyAnswered.Error(),
			Domain: constant.ServiceName,
			Metadata: map[string]string{
				"timestamp": nowRFC3339Nano(),
			},
		},
	)

	return ErrorInviteAlreadyAnswered.Raise(cause, st)
}
