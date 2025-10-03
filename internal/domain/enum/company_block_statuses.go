package enum

import "fmt"

const (
	DistributorBlockStatusActive   = "active"    // Distributor is active
	DistributorBlockStatusCanceled = "cancelled" // Distributor is canceled
)

var distributorBlockStatuses = []string{
	DistributorBlockStatusActive,
	DistributorBlockStatusCanceled,
}

var ErrorDistributorBlockStatusNotSupported = fmt.Errorf("block distributor status must be one of: %s", GetAllDistributorBlockStatuses())

func CheckDistributorBlockStatus(status string) error {
	for _, s := range distributorBlockStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", status, ErrorDistributorBlockStatusNotSupported)
}

func GetAllDistributorBlockStatuses() []string {
	return distributorBlockStatuses
}
