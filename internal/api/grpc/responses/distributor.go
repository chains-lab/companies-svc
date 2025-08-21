package responses

import (
	pagProto "github.com/chains-lab/distributors-proto/gen/go/common/pagination"
	disProto "github.com/chains-lab/distributors-proto/gen/go/svc/distributor"
	"github.com/chains-lab/distributors-svc/internal/app/models"
	"github.com/chains-lab/pagi"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Distributor(distributor models.Distributor) *disProto.Distributor {
	return &disProto.Distributor{
		Id:        distributor.ID.String(),
		Name:      distributor.Name,
		Status:    distributor.Status,
		CreatedAt: timestamppb.New(distributor.CreatedAt),
		UpdatedAt: timestamppb.New(distributor.UpdatedAt),
	}
}

func DistributorsList(distributors []models.Distributor, pag pagi.Response) *disProto.DistributorsList {
	list := make([]*disProto.Distributor, len(distributors))
	for i, distributor := range distributors {
		list[i] = Distributor(distributor)
	}

	return &disProto.DistributorsList{
		Distributors: list,
		Pagination: &pagProto.Response{
			Page:  pag.Page,
			Size:  pag.Size,
			Total: pag.Total,
		},
	}
}

func Block(block models.Block) *disProto.DistributorBlock {
	resp := &disProto.DistributorBlock{
		Id:            block.ID.String(),
		DistributorId: block.DistributorID.String(),
		InitiatorId:   block.InitiatorID.String(),
		Reason:        block.Reason,
		Status:        block.Status,
		BlockedAt:     timestamppb.New(block.BlockedAt),
	}

	if block.CanceledAt != nil {
		resp.CanceledAt = timestamppb.New(*block.CanceledAt)
	}

	return resp
}

func BlocksList(blocks []models.Block, pag pagi.Response) *disProto.DistributorBlocksList {
	resp := &disProto.DistributorBlocksList{
		Blocks: make([]*disProto.DistributorBlock, len(blocks)),
		Pagination: &pagProto.Response{
			Page:  pag.Page,
			Size:  pag.Size,
			Total: pag.Total,
		},
	}

	for i, block := range blocks {
		resp.Blocks[i] = Block(block)
	}

	return resp
}
