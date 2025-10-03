package enum

import "fmt"

const (
	DistributorStatusActive   = "active"
	DistributorStatusInactive = "inactive"
	DistributorStatusBlocked  = "blocked"
)

var distributorStatuses = []string{
	DistributorStatusActive,
	DistributorStatusInactive,
	DistributorStatusBlocked,
}

var ErrorDistributorStatusNotSupported = fmt.Errorf("distributor status not supported mus be one of: %v", GetAllDistributorStatuses())

func CheckDistributorStatus(status string) error {
	for _, s := range distributorStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", status, ErrorDistributorStatusNotSupported)
}

func GetAllDistributorStatuses() []string {
	return distributorStatuses
}
