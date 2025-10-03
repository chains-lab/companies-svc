package data

import (
	"context"
	"time"

	"github.com/chains-lab/companies-svc/internal/data/pgdb"
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/domain/service/company"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (d *Database) CreateCompanyBlock(ctx context.Context, input models.CompanyBlock) error {
	schema := blockModelToSchema(input)

	return d.sql.blockages.New().Insert(ctx, schema)
}

func (d *Database) GetCompanyBlockByID(ctx context.Context, ID uuid.UUID) (models.CompanyBlock, error) {
	schema, err := d.sql.blockages.New().FilterID(ID).Get(ctx)
	if err != nil {
		return models.CompanyBlock{}, err
	}

	return companyBlockSchemaToModel(schema), nil
}

func (d *Database) GetActiveCompanyBlock(ctx context.Context, companyID uuid.UUID) (models.CompanyBlock, error) {
	schema, err := d.sql.blockages.New().FiltercompanyID(companyID).FilterStatus(enum.DistributorStatusActive).Get(ctx)
	if err != nil {
		return models.CompanyBlock{}, err
	}

	return companyBlockSchemaToModel(schema), nil
}

func (d *Database) FilterCompanyBlocks(ctx context.Context, filters company.FilterBlockages, page, size uint64) (models.CompanyBlockCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := d.sql.blockages.New()

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
		return models.CompanyBlockCollection{}, err
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return models.CompanyBlockCollection{}, err
	}

	collection := make([]models.CompanyBlock, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, companyBlockSchemaToModel(row))
	}

	return models.CompanyBlockCollection{
		Data:  collection,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (d *Database) CancelActiveCompanyBlock(ctx context.Context, companyID uuid.UUID, canceledAt time.Time) error {
	return d.sql.blockages.New().
		FiltercompanyID(companyID).
		FilterStatus(enum.DistributorStatusActive).
		UpdateStatus(enum.DistributorBlockStatusActive).
		UpdateCanceledAt(canceledAt).
		Update(ctx)
}

func blockModelToSchema(m models.CompanyBlock) pgdb.CompanyBlock {
	block := pgdb.CompanyBlock{
		ID:          m.ID,
		CompanyID:   m.CompanyID,
		InitiatorID: m.InitiatorID,
		Reason:      m.Reason,
		Status:      m.Status,
		BlockedAt:   m.BlockedAt,
	}
	if m.CanceledAt != nil {
		block.CanceledAt = m.CanceledAt
	}

	return block
}

func companyBlockSchemaToModel(s pgdb.CompanyBlock) models.CompanyBlock {
	block := models.CompanyBlock{
		ID:          s.ID,
		CompanyID:   s.CompanyID,
		InitiatorID: s.InitiatorID,
		Reason:      s.Reason,
		Status:      s.Status,
		BlockedAt:   s.BlockedAt,
	}
	if s.CanceledAt != nil {
		block.CanceledAt = s.CanceledAt
	}

	return block
}
