package responses

import (
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/resources"
)

func CompanyBlock(m models.CompanyBlock) resources.CompanyBlock {
	resp := resources.CompanyBlock{
		Data: resources.CompanyBlockData{
			Id:   m.ID,
			Type: resources.CompanyBlockType,
			Attributes: resources.CompanyBlockAttributes{
				CompanyId:   m.CompanyID,
				InitiatorId: m.InitiatorID,
				Reason:      m.Reason,
				Status:      m.Status,
				BlockedAt:   m.BlockedAt,
			},
		},
	}

	if m.CanceledAt != nil {
		resp.Data.Attributes.CancelledAt = m.CanceledAt
	}

	return resp
}

func CompanyBlockCollection(ms models.CompanyBlockCollection) resources.CompanyBlocksCollection {
	items := make([]resources.CompanyBlockData, 0, len(ms.Data))
	for _, m := range ms.Data {
		items = append(items, CompanyBlock(m).Data)
	}

	return resources.CompanyBlocksCollection{
		Data: items,
		Links: resources.PaginationData{
			PageNumber: int64(ms.Page),
			PageSize:   int64(ms.Size),
			TotalItems: int64(ms.Total),
		},
	}
}
