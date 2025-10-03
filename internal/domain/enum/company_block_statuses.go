package enum

import "fmt"

const (
	CompanyBlockStatusActive   = "active"    // Company is active
	CompanyBlockStatusCanceled = "cancelled" // Company is canceled
)

var companyBlockStatuses = []string{
	CompanyBlockStatusActive,
	CompanyBlockStatusCanceled,
}

var ErrorCompanyBlockStatusNotSupported = fmt.Errorf("block company status must be one of: %s", GetAllCompanyBlockStatuses())

func CheckCompanyBlockStatus(status string) error {
	for _, s := range companyBlockStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", status, ErrorCompanyBlockStatusNotSupported)
}

func GetAllCompanyBlockStatuses() []string {
	return companyBlockStatuses
}
