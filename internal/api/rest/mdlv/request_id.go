package mdlv

import (
	"context"
	"net/http"

	"github.com/chains-lab/distributors-svc/internal/api/rest/meta"
	"github.com/google/uuid"
)

type ctxKey string

const (
	RequestIDHeader = "X-Request-ID"
)

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, meta.RequestIDCtxKey, id)
}

func FromRequestID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(meta.RequestIDCtxKey).(string)
	return v, ok
}

func RequestIDMdl() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get(RequestIDHeader)
			if reqID == "" {
				reqID = uuid.NewString()
			}
			ctx := WithRequestID(r.Context(), reqID)

			w.Header().Set(RequestIDHeader, reqID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
