package handlers

import (
	"context"

	"github.com/chains-lab/distributors-svc/pkg/logger"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ctxKey int

const (
	LogCtxKey ctxKey = iota
	RequestIDCtxKey
	UserCtxKey
)

func RequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	requestID, ok := ctx.Value(RequestIDCtxKey).(uuid.UUID)
	if !ok {
		return ""
	}

	return requestID.String()
}

type UserData struct {
	ID        uuid.UUID `json:"sub,omitempty"`
	SessionID uuid.UUID `json:"session_id,omitempty"`
	Verified  bool      `json:"verified,omitempty"`
	Role      string    `json:"role,omitempty"`
}

func User(ctx context.Context) (UserData, error) {
	if ctx == nil {
		return UserData{}, status.Error(codes.Internal, "internal server error")
	}

	userData, ok := ctx.Value(UserCtxKey).(UserData)
	if !ok {
		return UserData{}, status.Error(codes.Unauthenticated, "missing metadata in request")
	}

	return userData, nil
}

func Log(ctx context.Context) logger.Logger {
	entry, ok := ctx.Value(LogCtxKey).(logger.Logger)
	if !ok {
		logrus.Info("no logger in context")

		entry = logger.NewWithBase(logrus.New())
	}

	requestID := RequestID(ctx)

	if requestID != "" {
		return entry
	}

	return entry.WithField("request_id", requestID)
}
