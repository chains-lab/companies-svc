package responses

import (
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/distributors-svc/resources"
	"github.com/chains-lab/pagi"
)

func Distributor(m models.Distributor) resources.Distributor {
	resp := resources.Distributor{
		Data: resources.DistributorData{
			Id:   m.ID.String(),
			Type: resources.DistributorType,
			Attributes: resources.DistributorAttributes{
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

func DistributorCollection(ms []models.Distributor, pag pagi.Response) resources.DistributorsCollection {
	items := make([]resources.DistributorData, 0, len(ms))
	for _, m := range ms {
		items = append(items, Distributor(m).Data)
	}

	return resources.DistributorsCollection{
		Data: items,
		Links: resources.PaginationData{
			PageNumber: int64(pag.Page),
			PageSize:   int64(pag.Size),
			TotalItems: int64(pag.Total),
		},
	}
}
