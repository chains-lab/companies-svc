package responses

import (
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/distributors-svc/resources"
)

func DistributorBlock(m models.DistributorBlock) resources.DistributorBlock {
	resp := resources.DistributorBlock{
		Data: resources.DistributorBlockData{
			Id:   m.ID,
			Type: resources.DistributorBlockType,
			Attributes: resources.DistributorBlockAttributes{
				DistributorId: m.DistributorID,
				InitiatorId:   m.InitiatorID,
				Reason:        m.Reason,
				Status:        m.Status,
				BlockedAt:     m.BlockedAt,
			},
		},
	}

	if m.CanceledAt != nil {
		resp.Data.Attributes.CancelledAt = m.CanceledAt
	}

	return resp
}

func DistributorBlockCollection(ms models.DistributorBlockCollection) resources.DistributorBlocksCollection {
	items := make([]resources.DistributorBlockData, 0, len(ms.Data))
	for _, m := range ms.Data {
		items = append(items, DistributorBlock(m).Data)
	}

	return resources.DistributorBlocksCollection{
		Data: items,
		Links: resources.PaginationData{
			PageNumber: int64(ms.Page),
			PageSize:   int64(ms.Size),
			TotalItems: int64(ms.Total),
		},
	}
}
