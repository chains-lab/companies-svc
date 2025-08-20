package enum

import "fmt"

const (
	DistributorStatusActive   = "active"
	DistributorStatusInactive = "inactive"
	DistributorStatusSuspend  = "suspend"
)

var distributorStatuses = []string{
	DistributorStatusActive,
	DistributorStatusInactive,
	DistributorStatusSuspend,
}

var ErrorDistributorStatusNotSupported = fmt.Errorf("distributor status not supported mus be one of: %v", GetAllDistributorStatuses())

func ParseDistributorStatus(status string) (string, error) {
	for _, s := range distributorStatuses {
		if s == status {
			return s, nil
		}
	}

	return "", fmt.Errorf("'%s', %w", status, ErrorDistributorStatusNotSupported)
}

func GetAllDistributorStatuses() []string {
	return distributorStatuses
}
