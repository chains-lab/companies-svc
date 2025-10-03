package responses

import (
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/distributors-svc/resources"
)

func Distributor(m models.Distributor) resources.Distributor {
	resp := resources.Distributor{
		Data: resources.DistributorData{
			Id:   m.ID,
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

func DistributorCollection(ms models.DistributorCollection) resources.DistributorsCollection {
	items := make([]resources.DistributorData, 0, len(ms.Data))
	for _, m := range ms.Data {
		items = append(items, Distributor(m).Data)
	}

	return resources.DistributorsCollection{
		Data: items,
		Links: resources.PaginationData{
			PageNumber: int64(ms.Page),
			PageSize:   int64(ms.Size),
			TotalItems: int64(ms.Total),
		},
	}
}
