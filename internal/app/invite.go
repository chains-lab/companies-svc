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
	"github.com/chains-lab/distributors-svc/internal/problems"
	"github.com/chains-lab/pagi"
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

	OrderByCreatedAt(asc bool) dbx.InviteQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.InviteQ
}

func (a App) GetInvite(ctx context.Context, id uuid.UUID) (models.Invite, error) {
	invite, err := a.invite.New().FilterID(id).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Invite{}, problems.RaiseInviteNotFound(ctx, fmt.Errorf("invite not found: %w", err), id)
		default:
			return models.Invite{}, problems.RaiseInternal(ctx, fmt.Errorf("getting invite: %w", err))
		}
	}

	return models.Invite{
		ID:            invite.ID,
		DistributorID: invite.DistributorID,
		UserID:        invite.UserID,
		InvitedBy:     invite.InvitedBy,
		Role:          invite.Role,
		Status:        invite.Status,
		CreatedAt:     invite.CreatedAt,
		AnsweredAt:    invite.AnsweredAt,
	}, nil
}

func (a App) SendInvite(
	ctx context.Context,
	initiatorID uuid.UUID,
	userID uuid.UUID,
	distributorID uuid.UUID,
	role string,
) (models.Invite, error) {
	_, err := a.employee.New().FilterUserID(userID).Get(ctx)
	if err == nil || errors.Is(err, sql.ErrNoRows) {
		return models.Invite{}, problems.RaiseCantSendInviteForCurrentEmployee(
			ctx,
			fmt.Errorf("user already has an employee record: %s", userID),
			userID,
		)
	} else if err != nil {
		return models.Invite{}, problems.RaiseInternal(ctx, fmt.Errorf("getting user employee record: %w", err))
	}

	initiator, err := a.CompareEmployeesRole(ctx, initiatorID, distributorID, enum.EmployeeRoleAdmin)
	if err != nil {
		return models.Invite{}, err
	}

	access, err := enum.ComparisonEmployeeRoles(initiator.Role, role)
	if err != nil {
		return models.Invite{}, err
	}
	if access != 1 {
		return models.Invite{}, problems.RaiseInitiatorEmployeeHaveNotEnoughPermissions(
			ctx,
			fmt.Errorf("initiator have not enough rights"),
			initiatorID,
			distributorID,
		)
	}

	invite := dbx.Invite{
		ID:            uuid.New(),
		DistributorID: distributorID,
		UserID:        userID,
		InvitedBy:     initiatorID,
		Role:          role,
		Status:        enum.InviteStatusSent,
		CreatedAt:     time.Now().UTC(),
	}

	err = a.invite.New().Insert(ctx, invite)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Invite{}, problems.RaiseUserHaveAlreadyInviteForInitiatorDistributor(
				ctx,
				fmt.Errorf("user %s already has an invite for this distributor %s: %w", userID, distributorID, err),
				distributorID,
			)
		default:
			return models.Invite{}, problems.RaiseInternal(ctx, fmt.Errorf("inserting invite: %w", err))
		}
	}

	return models.Invite{}, nil
}

func (a App) SelectInvites(
	ctx context.Context,
	filters map[string]any,
	ascend bool,
	pag pagi.Request,
) ([]models.Invite, pagi.Response, error) {
	query := a.invite.New()

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

	limit, offset := pagi.CalculateLimitOffset(pag)

	count, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, problems.RaiseInternal(ctx, fmt.Errorf("counting invites: %w", err))
	}

	invites, err := query.Page(limit, offset).OrderByCreatedAt(ascend).Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, problems.RaiseInternal(ctx, fmt.Errorf("selecting invites: %w", err))
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
			CreatedAt:     i.CreatedAt,
		}
		if i.AnsweredAt != nil {
			element.AnsweredAt = i.AnsweredAt
		}
		res = append(res, element)
	}

	return res, pagi.Response{
		Page:  pag.Page,
		Size:  limit,
		Total: count,
	}, nil
}

func (a App) WithdrawInvite(
	ctx context.Context,
	initiatorID uuid.UUID,
	inviteID uuid.UUID,
) (models.Invite, error) {
	invite, err := a.invite.New().FilterID(inviteID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.Invite{}, problems.RaiseInviteNotFound(ctx, fmt.Errorf("invite not found: %w", err), inviteID)
		default:
			return models.Invite{}, problems.RaiseInternal(ctx, fmt.Errorf("getting invite: %w", err))
		}
	}

	_, err = a.CompareEmployeesRole(ctx, initiatorID, invite.DistributorID, enum.EmployeeRoleAdmin)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.Status != enum.InviteStatusSent {
		return models.Invite{}, problems.RaiseInviteIsNotActive(ctx, fmt.Errorf("invite already answered: %s", inviteID), inviteID)
	}

	err = a.invite.New().FilterID(inviteID).Update(ctx, map[string]interface{}{
		"status": enum.InviteStatusWithdrawn,
	})
	if err != nil {
		return models.Invite{}, problems.RaiseInternal(ctx, fmt.Errorf("updating invite: %w", err))
	}

	return models.Invite{
		ID:            invite.ID,
		DistributorID: invite.DistributorID,
		UserID:        invite.UserID,
		InvitedBy:     invite.InvitedBy,
		Role:          invite.Role,
		Status:        enum.InviteStatusWithdrawn,
		CreatedAt:     invite.CreatedAt,
	}, nil
}

func (a App) AcceptInvite(
	ctx context.Context,
	initiatorID uuid.UUID,
	inviteID uuid.UUID,
) (models.Invite, error) {
	invite, err := a.invite.New().FilterID(inviteID).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Invite{}, problems.RaiseInviteNotFound(ctx, fmt.Errorf("invite not found: %w", err), inviteID)
		}
		return models.Invite{}, problems.RaiseInternal(ctx, err)
	}

	if invite.Status != enum.InviteStatusSent {
		return models.Invite{}, problems.RaiseInviteIsNotActive(ctx, fmt.Errorf("invite already answered: %s", inviteID), inviteID)
	}

	if initiatorID != invite.UserID {
		return models.Invite{}, problems.RaiseInviteIsNotForInitiator(
			ctx,
			fmt.Errorf("invite is not for initiator: %s", inviteID),
			inviteID,
		)
	}

	_, err = a.employee.New().FilterUserID(invite.UserID).Get(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.Invite{}, problems.RaiseInternal(ctx, err)
	}
	if err == nil {
		return models.Invite{}, problems.RaiseInitiatorIsAlreadyEmployee(
			ctx,
			fmt.Errorf("initiator is already an employee: %s", initiatorID),
			initiatorID,
		)
	}

	invite.Status = enum.InviteStatusAccepted

	now := time.Now().UTC()

	trErr := a.employee.Transaction(func(ctx context.Context) error {
		err = a.invite.New().Update(ctx, map[string]any{
			"status":      invite.Status,
			"answered_at": now,
		})
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return problems.RaiseInviteNotFound(ctx, fmt.Errorf("invite not found: %w", err), inviteID)
			default:
				return problems.RaiseInternal(ctx, fmt.Errorf("updating invite: %w", err))
			}
		}

		err = a.employee.New().Insert(ctx, dbx.Employee{
			DistributorID: invite.DistributorID,
			UserID:        invite.UserID,
			Role:          invite.Role,
			UpdatedAt:     now,
			CreatedAt:     now,
		})
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				return problems.RaiseInviteNotFound(ctx, fmt.Errorf("invite not found: %w", err), inviteID)
			default:
				return problems.RaiseInternal(ctx, fmt.Errorf("inserting employee: %w", err))
			}
		}

		return nil
	})
	if trErr != nil {
		return models.Invite{}, trErr
	}

	return models.Invite{
		ID:            invite.ID,
		DistributorID: invite.DistributorID,
		UserID:        invite.UserID,
		InvitedBy:     invite.InvitedBy,
		Role:          invite.Role,
		Status:        invite.Status,
		CreatedAt:     invite.CreatedAt,
		AnsweredAt:    &now,
	}, nil
}

func (a App) RejectInvite(
	ctx context.Context,
	initiatorID uuid.UUID,
	inviteID uuid.UUID,
) (models.Invite, error) {
	invite, err := a.invite.New().FilterID(inviteID).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Invite{}, problems.RaiseInviteNotFound(ctx, fmt.Errorf("invite not found: %w", err), inviteID)
		}
		return models.Invite{}, problems.RaiseInternal(ctx, err)
	}

	if invite.Status != enum.InviteStatusSent {
		return models.Invite{}, problems.RaiseInviteIsNotActive(ctx, fmt.Errorf("invite already answered: %s", inviteID), inviteID)
	}

	if initiatorID != invite.UserID {
		return models.Invite{}, problems.RaiseInviteIsNotForInitiator(ctx, fmt.Errorf("invite is not for initiator: %s", inviteID), inviteID)
	}

	now := time.Now().UTC()

	err = a.invite.New().Update(ctx, map[string]any{
		"status":      enum.InviteStatusRejected,
		"answered_at": now,
	})
	if err != nil {
		return models.Invite{}, problems.RaiseInternal(ctx, fmt.Errorf("updating invite: %w", err))
	}

	return models.Invite{
		ID:            invite.ID,
		DistributorID: invite.DistributorID,
		UserID:        invite.UserID,
		InvitedBy:     invite.InvitedBy,
		Role:          invite.Role,
		Status:        enum.InviteStatusRejected,
		CreatedAt:     invite.CreatedAt,
		AnsweredAt:    &now,
	}, nil
}
