package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/internal/domain/service/company"
	"github.com/chains-lab/restkit/pagi"
	"github.com/google/uuid"
)

func (d *Database) CreateCompany(ctx context.Context, model models.Company) (models.Company, error) {
	schema := companyModelToSchema(model)

	err := d.sql.companies.New().Insert(ctx, schema)
	if err != nil {
		return models.Company{}, err
	}

	return companiesSchemaToModel(schema), nil
}

func (d *Database) GetCompanyByID(ctx context.Context, ID uuid.UUID) (models.Company, error) {
	schema, err := d.sql.companies.New().FilterID(ID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Company{}, nil
	case err != nil:
		return models.Company{}, err
	}

	return companiesSchemaToModel(schema), nil
}

func (d *Database) FilterCompanies(
	ctx context.Context,
	filters company.FiltersParams,
	page, size uint64,
) (models.CompaniesCollection, error) {
	limit, offset := pagi.PagConvert(page, size)

	query := d.sql.companies.New()

	if filters.Statuses != nil {
		query = query.FilterStatus(filters.Statuses...)
	}
	if filters.Name != nil {
		query = query.FilterLikeName(*filters.Name)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return models.CompaniesCollection{}, err
	}

	rows, err := query.Page(limit, offset).Select(ctx)
	if err != nil {
		return models.CompaniesCollection{}, err
	}

	collection := make([]models.Company, 0, len(rows))
	for _, r := range rows {
		collection = append(collection, companiesSchemaToModel(r))
	}

	return models.CompaniesCollection{
		Data:  collection,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

func (d *Database) UpdateCompany(
	ctx context.Context,
	ID uuid.UUID,
	params company.UpdateParams,
	updatedAt time.Time,
) error {
	q := d.sql.companies.New().FilterID(ID)

	if params == (company.UpdateParams{}) {
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

func (d *Database) UpdateCompaniesStatus(
	ctx context.Context,
	ID uuid.UUID,
	status string,
	updatedAt time.Time,
) error {
	err := d.sql.companies.New().FilterID(ID).UpdateStatus(status).Update(ctx, updatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) DeleteCompany(ctx context.Context, ID uuid.UUID) error {
	err := d.sql.companies.New().FilterID(ID).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
