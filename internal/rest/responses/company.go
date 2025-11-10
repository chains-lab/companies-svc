package responses

import (
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/resources"
)

func Company(m models.Company) resources.Company {
	resp := resources.Company{
		Data: resources.CompanyData{
			Id:   m.ID,
			Type: resources.CompanyType,
			Attributes: resources.CompanyAttributes{
				Icon:      m.Icon,
				Name:      m.Name,
				Status:    m.Status,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			},
		},
	}

	return resp
}

func CompanyCollection(ms models.CompaniesCollection) resources.CompaniesCollection {
	items := make([]resources.CompanyData, 0, len(ms.Data))
	for _, m := range ms.Data {
		items = append(items, Company(m).Data)
	}

	return resources.CompaniesCollection{
		Data: items,
		Links: resources.PaginationData{
			PageNumber: int64(ms.Page),
			PageSize:   int64(ms.Size),
			TotalItems: int64(ms.Total),
		},
	}
}
