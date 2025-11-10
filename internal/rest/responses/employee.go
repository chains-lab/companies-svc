package responses

import (
	"github.com/chains-lab/companies-svc/internal/domain/models"
	"github.com/chains-lab/companies-svc/resources"
)

func Employee(m models.Employee) resources.Employee {
	resp := resources.Employee{
		Data: resources.EmployeeData{
			Id:   m.UserID,
			Type: resources.EmployeeType,
			Attributes: resources.EmployeeAttributes{
				CompanyId: m.CompanyID,
				Role:      m.Role,
				Position:  m.Position,
				Label:     m.Label,
				CreatedAt: m.CreatedAt,
				UpdatedAt: m.UpdatedAt,
			},
		},
	}

	return resp
}

func EmployeeCollection(ms models.EmployeesCollection) resources.EmployeesCollection {
	items := make([]resources.EmployeeData, 0, len(ms.Data))
	for _, m := range ms.Data {
		items = append(items, Employee(m).Data)
	}

	return resources.EmployeesCollection{
		Data: items,
		Links: resources.PaginationData{
			PageNumber: int64(ms.Page),
			PageSize:   int64(ms.Size),
			TotalItems: int64(ms.Total),
		},
	}
}
