package enum

import "fmt"

const (
	CompanyStatusActive   = "active"
	CompanyStatusInactive = "inactive"
	CompanyStatusBlocked  = "blocked"
)

var companyStatuses = []string{
	CompanyStatusActive,
	CompanyStatusInactive,
	CompanyStatusBlocked,
}

var ErrorCompanyStatusNotSupported = fmt.Errorf("company status not supported mus be one of: %v", GetAllCompanyStatuses())

func CheckCompanyStatus(status string) error {
	for _, s := range companyStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", status, ErrorCompanyStatusNotSupported)
}

func GetAllCompanyStatuses() []string {
	return companyStatuses
}
