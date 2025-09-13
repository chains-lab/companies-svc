package jwtmanager

import (
	"time"

	"github.com/chains-lab/distributors-svc/internal/config"
	"github.com/chains-lab/enum"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Manager подписывает/проверяет инвайт-JWT.
type Manager struct {
	iss string
	sk  string
}

type InviteData struct {
	JTI           uuid.UUID
	DistributorID uuid.UUID
	Role          string
	ExpiresAt     time.Time
	Issuer        string
}

// наши claims внутри JWT
type inviteClaims struct {
	DistributorID uuid.UUID `json:"distributor_id"`
	Role          string    `json:"role"`
	jwt.RegisteredClaims
}

func NewManager(cfg config.Config) Manager {
	return Manager{
		iss: enum.CitiesSVC,
		sk:  cfg.JWT.Invites.SecretKey,
	}
}
