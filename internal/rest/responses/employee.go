package responses

import (
	"github.com/chains-lab/distributors-svc/internal/domain/models"
	"github.com/chains-lab/distributors-svc/resources"
	"github.com/chains-lab/pagi"
)

func Employee(m models.Employee) resources.Employee {
	resp := resources.Employee{
		Data: resources.EmployeeData{
			Id:   m.UserID.String(),
			Type: resources.EmployeeType,
			Attributes: resources.EmployeeAttributes{
				DistributorId: m.DistributorID.String(),
				Role:          m.Role,
				CreatedAt:     m.CreatedAt,
				UpdatedAt:     m.UpdatedAt,
			},
		},
	}

	return resp
}

func EmployeeCollection(ms []models.Employee, pag pagi.Response) resources.EmployeesCollection {
	items := make([]resources.EmployeeData, 0, len(ms))
	for _, m := range ms {
		items = append(items, Employee(m).Data)
	}

	return resources.EmployeesCollection{
		Data: items,
		Links: resources.PaginationData{
			PageNumber: int64(pag.Page),
			PageSize:   int64(pag.Size),
			TotalItems: int64(pag.Total),
		},
	}
}
