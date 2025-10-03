package data

import (
	"context"
	"time"

	"github.com/chains-lab/distributors-svc/internal/data/pgdb"
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/distributors-svc/internal/domain/service/distributor"
	"github.com/chains-lab/pagi"
	"github.com/google/uuid"
)

func (d *Database) CreateDistributor(ctx context.Context, model models.Distributor) (models.Distributor, error) {
	schema := distributorModelToSchema(model)

	err := d.sql.distributors.New().Insert(ctx, schema)
	if err != nil {
		return models.Distributor{}, err
	}

	return distributorSchemaToModel(schema), nil
}

func (d *Database) GetDistributorByID(ctx context.Context, ID uuid.UUID) (models.Distributor, error) {
	schema, err := d.sql.distributors.New().FilterID(ID).Get(ctx)
	if err != nil {
		return models.Distributor{}, err
	}

	return distributorSchemaToModel(schema), nil
}

func (d *Database) FilterDistributors(
	ctx context.Context,
	filters distributor.Filters,
	page, size uint64,
) (models.DistributorCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := d.sql.distributors.New()

	if filters.Statuses != nil {
		query = query.FilterStatus(filters.Statuses...)
	}
	if filters.Name != nil {
		query = query.FilterLikeName(*filters.Name)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return models.DistributorCollection{}, err
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return models.DistributorCollection{}, err
	}

	collection := make([]models.Distributor, 0, len(rows))
	for _, r := range rows {
		collection = append(collection, distributorSchemaToModel(r))
	}

	return models.DistributorCollection{
		Data:  collection,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (d *Database) UpdateDistributor(
	ctx context.Context,
	ID uuid.UUID,
	params distributor.UpdateParams,
	updatedAt time.Time,
) error {
	q := d.sql.distributors.New().FilterID(ID)

	if params == (distributor.UpdateParams{}) {
		return nil
	}

	if params.Name != nil {
		q = q.UpdateName(*params.Name)
	}
	if params.Icon != nil {
		q = q.UpdateIcon(*params.Icon)
	}

	return q.Update(ctx, updatedAt)
}

func (d *Database) UpdateDistributorStatus(
	ctx context.Context,
	ID uuid.UUID,
	status string,
	updatedAt time.Time,
) error {
	return d.sql.distributors.New().FilterID(ID).UpdateStatus(status).Update(ctx, updatedAt)
}

func distributorModelToSchema(model models.Distributor) pgdb.Distributor {
	return pgdb.Distributor{
		ID:        model.ID,
		Name:      model.Name,
		Icon:      model.Icon,
		Status:    model.Status,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func distributorSchemaToModel(schema pgdb.Distributor) models.Distributor {
	return models.Distributor{
		ID:        schema.ID,
		Name:      schema.Name,
		Icon:      schema.Icon,
		Status:    schema.Status,
		CreatedAt: schema.CreatedAt,
		UpdatedAt: schema.UpdatedAt,
	}
}
