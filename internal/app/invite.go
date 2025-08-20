package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/internal/config/constant/enum"
	"github.com/chains-lab/distributors-svc/internal/dbx"
	"github.com/chains-lab/distributors-svc/internal/errx"
	"github.com/chains-lab/distributors-svc/pkg/pagination"
	"github.com/google/uuid"
)

type inviteQ interface {
	New() dbx.InviteQ

	Insert(ctx context.Context, input dbx.Invite) error
	Get(ctx context.Context) (dbx.Invite, error)
	Select(ctx context.Context) ([]dbx.Invite, error)
	Update(ctx context.Context, input map[string]any) error
	Delete(ctx context.Context) error

	FilterID(id uuid.UUID) dbx.InviteQ
	FilterDistributorID(distributorID uuid.UUID) dbx.InviteQ
	FilterUserID(userID uuid.UUID) dbx.InviteQ
	FilterInvitedBy(userID uuid.UUID) dbx.InviteQ
	FilterRole(role string) dbx.InviteQ
	FilterStatus(status string) dbx.InviteQ

	OrderByExpiresAt(asc bool) dbx.InviteQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.InviteQ
}

func (a App) GetInvite(ctx context.Context, id uuid.UUID) (models.Invite, error) {
	invite, err := a.invite.New().FilterID(id).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Invite{}, errx.RaiseInviteNotFound(ctx, fmt.Errorf("invite not found: %w", err), id)
		}
		return models.Invite{}, errx.RaiseInternal(ctx, err)
	}

	return models.Invite{
		ID:            invite.ID,
		DistributorID: invite.DistributorID,
		UserID:        invite.UserID,
		InvitedBy:     invite.InvitedBy,
		Role:          invite.Role,
		Status:        invite.Status,
		ExpiresAt:     invite.ExpiresAt,
		CreatedAt:     invite.CreatedAt,
		AnsweredAt:    invite.AnsweredAt,
	}, nil
}

func (a App) SendInvite(ctx context.Context, initiatorID, userID, cityID uuid.UUID, role string) (models.Invite, error) {
	_, err := a.employee.New().FilterUserID(userID).Get(ctx)
	if err == nil || errors.Is(err, sql.ErrNoRows) {
		return models.Invite{}, errx.UserAlreadyEmployee(
			ctx,
			fmt.Errorf("user already has an employee record: %s", userID),
			userID,
		)
	} else if err != nil {
		return models.Invite{}, errx.RaiseInternal(ctx, fmt.Errorf("getting user employee record: %w", err))
	}

	initiator, err := a.CompareEmployeesRole(ctx, initiatorID, cityID, enum.EmployeeRoleAdmin)
	if err != nil {
		return models.Invite{}, err
	}

	access, err := enum.ComparisonEmployeeRoles(initiator.Role, role)
	if err != nil {
		return models.Invite{}, err
	}
	if access != 1 {
		return models.Invite{}, errx.RaiseInitiatorEmployeeHaveNotEnoughPermissions(
			ctx,
			fmt.Errorf("initiator have not enough rights"),
			initiatorID,
			cityID,
		)
	}

	invite := dbx.Invite{
		ID:            uuid.New(),
		DistributorID: cityID,
		UserID:        userID,
		InvitedBy:     initiatorID,
		Role:          role,
		Status:        enum.InviteStatusSent,
		ExpiresAt:     time.Now().UTC().Add(24 * time.Hour), // 24 hours expiration
		CreatedAt:     time.Now().UTC(),
	}

	err = a.invite.New().Insert(ctx, invite)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Invite{}, err // TODO user already have invite
		default:
			return models.Invite{}, errx.RaiseInternal(ctx, err)
		}
	}

	return models.Invite{}, nil
}

func (a App) GetInvites(ctx context.Context, filters map[string]any, pag pagination.Request) ([]models.Invite, pagination.Response, error) {
	query := a.invite.New()

	if id, ok := filters["id"].(uuid.UUID); ok {
		query = query.FilterID(id)
	}
	if distributorID, ok := filters["distributor_id"].(uuid.UUID); ok {
		query = query.FilterDistributorID(distributorID)
	}
	if userID, ok := filters["user_id"].(uuid.UUID); ok {
		query = query.FilterUserID(userID)
	}
	if invitedBy, ok := filters["invited_by"].(uuid.UUID); ok {
		query = query.FilterInvitedBy(invitedBy)
	}
	if role, ok := filters["role"].(string); ok {
		query = query.FilterRole(role)
	}
	if status, ok := filters["status"].(string); ok {
		query = query.FilterStatus(status)
	}

	limit, offset := pagination.CalculateLimitOffset(pag)

	invites, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, pagination.Response{}, errx.RaiseInternal(ctx, err)
	}

	res := make([]models.Invite, 0, len(invites))

	for _, i := range invites {
		element := models.Invite{
			ID:            i.ID,
			DistributorID: i.DistributorID,
			UserID:        i.UserID,
			InvitedBy:     i.InvitedBy,
			Role:          i.Role,
			Status:        i.Status,
			ExpiresAt:     i.ExpiresAt,
			CreatedAt:     i.CreatedAt,
		}
		if i.AnsweredAt != nil {
			element.AnsweredAt = i.AnsweredAt
		}
		res = append(res, element)
	}

	return res, pagination.Response{
		Page:  pag.Page,
		Size:  limit,
		Total: count,
	}, nil
}

func (a App) AcceptInvite(ctx context.Context, inviteID uuid.UUID) error {
	invite, err := a.invite.New().FilterID(inviteID).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errx.RaiseInviteNotFound(ctx, fmt.Errorf("invite not found: %w", err), inviteID)
		}
		return errx.RaiseInternal(ctx, err)
	}

	if invite.Status != enum.InviteStatusSent {
		return errx.RaiseInviteAlreadyAccepted(ctx, fmt.Errorf("invite already answered: %s", inviteID))
	}

	_, err = a.employee.New().FilterUserID(invite.UserID).Get(ctx)
	if err == nil || errors.Is(err, sql.ErrNoRows) {
		return errx.RaiseInviteAlreadyAccepted(ctx, fmt.Errorf("user already has an employee record: %s", invite.UserID))
	} else if err != nil {
		return errx.RaiseInternal(ctx, err)
	}

	invite.Status = enum.InviteStatusAccepted
	invite.ExpiresAt = time.Time{} // Clear expiration date

	err = a.employee.Transaction(func(ctx context.Context) error {
		err = a.invite.New().Update(ctx, map[string]any{
			"status":     invite.Status,
			"expires_at": invite.ExpiresAt,
		})
		if err != nil {
			return errx.RaiseInternal(ctx, fmt.Errorf("updating invite: %w", err))
		}

		err = a.employee.New().Insert(ctx, dbx.Employee{
			DistributorID: invite.DistributorID,
			UserID:        invite.UserID,
			Role:          invite.Role,
			UpdatedAt:     time.Now().UTC(),
			CreatedAt:     time.Now().UTC(),
		})
		if err != nil {
			return errx.RaiseInternal(ctx, fmt.Errorf("inserting employee record: %w", err))
		}

		return nil
	})

	return nil
}

func (a App) RejectInvite(ctx context.Context, inviteID uuid.UUID) error {
	invite, err := a.invite.New().FilterID(inviteID).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errx.RaiseInviteNotFound(ctx, fmt.Errorf("invite not found: %w", err), inviteID)
		}
		return errx.RaiseInternal(ctx, err)
	}

	if invite.Status != enum.InviteStatusSent {
		return errx.RaiseInviteAlreadyAccepted(ctx, fmt.Errorf("invite already answered: %s", inviteID))
	}

	invite.Status = enum.InviteStatusRejected
	invite.ExpiresAt = time.Time{} // Clear expiration date

	err = a.invite.New().Update(ctx, map[string]any{
		"status":     invite.Status,
		"expires_at": invite.ExpiresAt,
	})
	if err != nil {
		return errx.RaiseInternal(ctx, fmt.Errorf("updating invite: %w", err))
	}

	return nil
}
