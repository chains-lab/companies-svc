package models

import (
	"time"

	"github.com/google/uuid"
)

type Invite struct {
	ID         uuid.UUID
	Status     string
	Role       string
	CompanyID  uuid.UUID
	UserID     *uuid.UUID
	AnsweredAt *time.Time
	Token      string
	ExpiresAt  time.Time
	CreatedAt  time.Time
}

func (i Invite) IsNil() bool {
	return i.ID == uuid.Nil
}

type InviteTokenData struct {
	InviteID  uuid.UUID
	CityID    uuid.UUID
	Role      string
	ExpiresAt time.Time
}

func (i InviteTokenData) IsNil() bool {
	return i == InviteTokenData{}
}
