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

	FilterID(ID uuid.UUID) dbx.InviteQ
	FilterDistributorID(distributorID ...uuid.UUID) dbx.InviteQ
	FilterUserID(userID ...uuid.UUID) dbx.InviteQ
	FilterInvitedBy(userID ...uuid.UUID) dbx.InviteQ
	FilterRole(role ...string) dbx.InviteQ
	FilterStatus(status ...string) dbx.InviteQ

	OrderByCreatedAt(asc bool) dbx.InviteQ

	Count(ctx context.Context) (uint64, error)
	Page(limit, offset uint64) dbx.InviteQ
}

func (a App) GetInvite(ctx context.Context, id uuid.UUID) (models.EmployeeInvite, error) {
	invite, err := a.invite.New().FilterID(id).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.EmployeeInvite{}, errx.InviteNotFound.Raise(
				fmt.Errorf("invite not found: %w", err),
			)
		default:
			return models.EmployeeInvite{}, errx.Internal.Raise(
				fmt.Errorf("getting invite: %w", err),
			)
		}
	}

	return models.EmployeeInvite{
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

func (a App) CreateEmployeeInvite(
	ctx context.Context,
	initiatorID uuid.UUID,
	userID uuid.UUID,
	distributorID uuid.UUID,
	role string,
) (models.EmployeeInvite, error) {
	_, err := a.employee.New().FilterUserID(userID).Get(ctx)
	if err == nil || errors.Is(err, sql.ErrNoRows) {
		return models.EmployeeInvite{}, errx.CantSendInviteForCurrentEmployee.Raise(
			fmt.Errorf("user already has an employee record: %s", userID),
		)
	} else if err != nil {
		return models.EmployeeInvite{}, errx.Internal.Raise(fmt.Errorf("getting user employee record: %w", err))
	}

	initiator, err := a.GetInitiatorEmployee(ctx, initiatorID)
	if err != nil {
		return models.EmployeeInvite{}, err
	}

	access, err := enum.ComparisonEmployeeRoles(initiator.Role, role)
	if err != nil {
		return models.EmployeeInvite{}, errx.EmployeeRoleNotSupported.Raise(
			fmt.Errorf("role %s not supported: %w", role, err),
		)
	}
	if access < 1 {
		return models.EmployeeInvite{}, errx.InitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator %s have not enough permissions to update distributor %s", initiatorID, distributorID),
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
			return models.EmployeeInvite{}, errx.UserHaveAlreadyInviteForInitiatorDistributor.Raise(
				fmt.Errorf("user %s already has an invite for this distributor %s: %w", userID, distributorID, err),
			)
		default:
			return models.EmployeeInvite{}, errx.Internal.Raise(fmt.Errorf("inserting invite: %w", err))
		}
	}

	return models.EmployeeInvite{}, nil
}

type SelectInvitesParams struct {
	Distributors []uuid.UUID
	ForUsers     []uuid.UUID
	Inviters     []uuid.UUID
	Roles        []string
	Statuses     []string
}

func (a App) SelectInvites(
	ctx context.Context,
	filters SelectInvitesParams,
	pag pagi.Request,
	sort []pagi.SortField,
) ([]models.EmployeeInvite, pagi.Response, error) {
	query := a.invite.New()

	if filters.Distributors != nil {
		query = query.FilterDistributorID(filters.Distributors...)
	}
	if filters.ForUsers != nil {
		query = query.FilterUserID(filters.ForUsers...)
	}
	if filters.Inviters != nil {
		query = query.FilterInvitedBy(filters.Inviters...)
	}
	if filters.Roles != nil {
		query = query.FilterRole(filters.Roles...)
	}
	if filters.Statuses != nil {
		query = query.FilterStatus(filters.Statuses...)
	}

	if pag.Page == 0 {
		pag.Page = 1
	}
	if pag.Size == 0 {
		pag.Size = 20
	}
	if pag.Size > 100 {
		pag.Size = 100
	}

	limit := pag.Size + 1
	offset := (pag.Page - 1) * pag.Size

	query = query.Page(limit, offset)

	for _, sort := range sort {
		ascend := sort.Ascend
		switch sort.Field {
		case "created_at":
			query = query.OrderByCreatedAt(ascend)
		default:

		}
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.Internal.Raise(
			fmt.Errorf("counting rows: %w", err),
		)
	}

	rows, err := query.Select(ctx)
	if err != nil {
		return nil, pagi.Response{}, errx.Internal.Raise(
			fmt.Errorf("selecting rows: %w", err),
		)
	}

	if len(rows) == int(limit) {
		rows = rows[:pag.Size]
	}

	res := make([]models.EmployeeInvite, 0, len(rows))

	for _, i := range rows {
		element := models.EmployeeInvite{
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
) (models.EmployeeInvite, error) {
	invite, err := a.invite.New().FilterID(inviteID).Get(ctx)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return models.EmployeeInvite{}, errx.InviteNotFound.Raise(
				fmt.Errorf("invite not found: %w", err),
			)
		default:
			return models.EmployeeInvite{}, errx.Internal.Raise(
				fmt.Errorf("getting invite: %w", err),
			)
		}
	}

	initiator, err := a.GetInitiatorEmployee(ctx, initiatorID)
	if err != nil {
		return models.EmployeeInvite{}, err
	}

	access, err := enum.ComparisonEmployeeRoles(initiator.Role, enum.EmployeeRoleAdmin)
	if err != nil {
		return models.EmployeeInvite{}, errx.EmployeeRoleNotSupported.Raise(
			fmt.Errorf("role %s not supported: %w", initiator.Role, err),
		)
	}
	if access < 0 {
		return models.EmployeeInvite{}, errx.InitiatorEmployeeHaveNotEnoughRights.Raise(
			fmt.Errorf("initiator %s have not enough permissions to withdraw invite %s", initiatorID, inviteID),
		)
	}

	if invite.Status != enum.InviteStatusSent {
		return models.EmployeeInvite{}, errx.InviteIsNotActive.Raise(
			fmt.Errorf("invite already answered: %s", inviteID),
		)
	}

	err = a.invite.New().FilterID(inviteID).Update(ctx, map[string]interface{}{
		"status": enum.InviteStatusWithdrawn,
	})
	if err != nil {
		return models.EmployeeInvite{}, errx.Internal.Raise(
			fmt.Errorf("updating invite: %w", err),
		)
	}

	return models.EmployeeInvite{
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
) (models.EmployeeInvite, error) {
	invite, err := a.invite.New().FilterID(inviteID).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.EmployeeInvite{}, errx.InviteNotFound.Raise(
				fmt.Errorf("invite not found: %w", err),
			)
		}
		return models.EmployeeInvite{}, errx.Internal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	if invite.Status != enum.InviteStatusSent {
		return models.EmployeeInvite{}, errx.InviteIsNotActive.Raise(
			fmt.Errorf("invite already answered: %s", inviteID),
		)
	}

	if initiatorID != invite.UserID {
		return models.EmployeeInvite{}, errx.InviteIsNotForInitiator.Raise(
			fmt.Errorf("invite is not for initiator: %s", inviteID),
		)
	}

	_, err = a.employee.New().FilterUserID(invite.UserID).Get(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.EmployeeInvite{}, errx.Internal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}
	if err == nil {
		return models.EmployeeInvite{}, errx.InitiatorIsAlreadyEmployee.Raise(
			fmt.Errorf("initiator is already an employee: %s", initiatorID),
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
				return errx.InviteNotFound.Raise(
					fmt.Errorf("invite not found: %w", err),
				)
			default:
				return errx.Internal.Raise(
					fmt.Errorf("updating invite: %w", err),
				)
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
				return errx.InviteNotFound.Raise(
					fmt.Errorf("invite not found: %w", err),
				)
			default:
				return errx.Internal.Raise(
					fmt.Errorf("inserting employee: %w", err),
				)
			}
		}

		return nil
	})
	if trErr != nil {
		return models.EmployeeInvite{}, trErr
	}

	return models.EmployeeInvite{
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
) (models.EmployeeInvite, error) {
	invite, err := a.invite.New().FilterID(inviteID).Get(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.EmployeeInvite{}, errx.InviteNotFound.Raise(
				fmt.Errorf("invite not found: %w", err),
			)
		}
		return models.EmployeeInvite{}, errx.Internal.Raise(
			fmt.Errorf("internal error: %w", err),
		)
	}

	if invite.Status != enum.InviteStatusSent {
		return models.EmployeeInvite{}, errx.InviteIsNotActive.Raise(
			fmt.Errorf("invite already answered: %s", inviteID),
		)
	}

	if initiatorID != invite.UserID {
		return models.EmployeeInvite{}, errx.InviteIsNotForInitiator.Raise(
			fmt.Errorf("invite is not for initiator: %s", inviteID),
		)
	}

	now := time.Now().UTC()

	err = a.invite.New().Update(ctx, map[string]any{
		"status":      enum.InviteStatusRejected,
		"answered_at": now,
	})
	if err != nil {
		return models.EmployeeInvite{}, errx.Internal.Raise(
			fmt.Errorf("updating invite: %w", err),
		)
	}

	return models.EmployeeInvite{
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
