package data

import (
	"context"
	"time"

	"github.com/chains-lab/distributors-svc/internal/data/pgdb"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/distributors-svc/internal/domain/service/distributor"
	"github.com/chains-lab/enum"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (d *Database) CreateDistributorBlock(ctx context.Context, input models.DistributorBlock) error {
	schema := blockModelToSchema(input)

	return d.sql.blockages.New().Insert(ctx, schema)
}

func (d *Database) GetDistributorBlockByID(ctx context.Context, ID uuid.UUID) (models.DistributorBlock, error) {
	schema, err := d.sql.blockages.New().FilterID(ID).Get(ctx)
	if err != nil {
		return models.DistributorBlock{}, err
	}

	return distributorBlockSchemaToModel(schema), nil
}

func (d *Database) GetActiveDistributorBlock(ctx context.Context, distributorID uuid.UUID) (models.DistributorBlock, error) {
	schema, err := d.sql.blockages.New().FilterDistributorID(distributorID).FilterStatus(enum.DistributorStatusActive).Get(ctx)
	if err != nil {
		return models.DistributorBlock{}, err
	}

	return distributorBlockSchemaToModel(schema), nil
}

func (d *Database) FilterDistributorBlocks(ctx context.Context, filters distributor.FilterBlockages, page, size uint64) (models.DistributorBlockCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := d.sql.blockages.New()

	if filters.Status != nil {
		query = query.FilterStatus(*filters.Status)
	}
	if filters.DistributorID != nil {
		query = query.FilterDistributorID(*filters.DistributorID)
	}
	if filters.InitiatorID != nil {
		query = query.FilterInitiatorID(*filters.InitiatorID)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return models.DistributorBlockCollection{}, err
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return models.DistributorBlockCollection{}, err
	}

	collection := make([]models.DistributorBlock, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, distributorBlockSchemaToModel(row))
	}

	return models.DistributorBlockCollection{
		Data:  collection,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (d *Database) CancelActiveDistributorBlock(ctx context.Context, distributorID uuid.UUID, canceledAt time.Time) error {
	return d.sql.blockages.New().
		FilterDistributorID(distributorID).
		FilterStatus(enum.DistributorStatusActive).
		UpdateStatus(enum.DistributorBlockStatusCanceled).
		UpdateCanceledAt(canceledAt).
		Update(ctx)
}

func blockModelToSchema(m models.DistributorBlock) pgdb.DistributorBlock {
	block := pgdb.DistributorBlock{
		ID:            m.ID,
		DistributorID: m.DistributorID,
		InitiatorID:   m.InitiatorID,
		Reason:        m.Reason,
		Status:        m.Status,
		BlockedAt:     m.BlockedAt,
	}
	if m.CanceledAt != nil {
		block.CanceledAt = m.CanceledAt
	}

	return block
}

func distributorBlockSchemaToModel(s pgdb.DistributorBlock) models.DistributorBlock {
	block := models.DistributorBlock{
		ID:            s.ID,
		DistributorID: s.DistributorID,
		InitiatorID:   s.InitiatorID,
		Reason:        s.Reason,
		Status:        s.Status,
		BlockedAt:     s.BlockedAt,
	}
	if s.CanceledAt != nil {
		block.CanceledAt = s.CanceledAt
	}

	return block
}
