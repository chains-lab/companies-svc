package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/enum"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/domain/services/block"
	"github.com/chains-lab/restkit/pagi"
	"github.com/google/uuid"
)

func (r *Repo) CreateCompanyBlock(ctx context.Context, input models.CompanyBlock) error {
	return r.sql.blockages.New().Insert(ctx, blockModelToSchema(input))
}

func (r *Repo) GetCompanyBlockByID(ctx context.Context, ID uuid.UUID) (models.CompanyBlock, error) {
	schema, err := r.sql.blockages.New().FilterID(ID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.CompanyBlock{}, nil
	case err != nil:
		return models.CompanyBlock{}, err
	}

	return companyBlockSchemaToModel(schema), nil
}

func (r *Repo) GetActiveCompanyBlock(ctx context.Context, companyID uuid.UUID) (models.CompanyBlock, error) {
	schema, err := r.sql.blockages.New().FiltercompanyID(companyID).FilterStatus(enum.CompanyStatusActive).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.CompanyBlock{}, nil
	case err != nil:
		return models.CompanyBlock{}, err
	}

	return companyBlockSchemaToModel(schema), nil
}

func (r *Repo) FilterCompanyBlocks(
	ctx context.Context,
	filters block.FilterParams,
	page, size uint64,
) (models.CompanyBlocksCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := r.sql.blockages.New()

	if filters.Status != nil {
		query = query.FilterStatus(*filters.Status)
	}
	if filters.CompanyID != nil {
		query = query.FiltercompanyID(*filters.CompanyID)
	}
	if filters.InitiatorID != nil {
		query = query.FilterInitiatorID(*filters.InitiatorID)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return models.CompanyBlocksCollection{}, err
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return models.CompanyBlocksCollection{}, err
	}

	collection := make([]models.CompanyBlock, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, companyBlockSchemaToModel(row))
	}

	return models.CompanyBlocksCollection{
		Data:  collection,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (r *Repo) CancelActiveCompanyBlock(ctx context.Context, companyID uuid.UUID, canceledAt time.Time) error {
	return r.sql.blockages.New().
		FiltercompanyID(companyID).
		FilterStatus(enum.CompanyStatusActive).
		UpdateStatus(enum.CompanyBlockStatusActive).
		UpdateCanceledAt(canceledAt).
		Update(ctx)
}
