package responses

import (
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/distributors-svc/resources"
	"github.com/chains-lab/pagi"
)

func DistributorBlock(m models.DistributorBlock) resources.DistributorBlock {
	resp := resources.DistributorBlock{
		Data: resources.DistributorBlockData{
			Id:   m.ID.String(),
			Type: resources.DistributorBlockType,
			Attributes: resources.DistributorBlockAttributes{
				DistributorId: m.DistributorID.String(),
				InitiatorId:   m.InitiatorID.String(),
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

func DistributorBlockCollection(ms []models.DistributorBlock, pag pagi.Response) resources.DistributorBlocksCollection {
	items := make([]resources.DistributorBlockData, 0, len(ms))
	for _, m := range ms {
		items = append(items, DistributorBlock(m).Data)
	}

	return resources.DistributorBlocksCollection{
		Data: items,
		Links: resources.PaginationData{
			PageNumber: int64(pag.Page),
			PageSize:   int64(pag.Size),
			TotalItems: int64(pag.Total),
		},
	}
}
