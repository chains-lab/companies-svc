package enum

import "fmt"

const (
	SuspendedDistributorStatusActive   = "active"    // Distributor is active
	SuspendedDistributorStatusCanceled = "cancelled" // Distributor is canceled
)

var suspendedDistributorStatuses = []string{
	SuspendedDistributorStatusActive,
	SuspendedDistributorStatusCanceled,
}

var ErrorSuspendedDistributorStatusNotSupported = fmt.Errorf("suspended distributor status must be one of: %s", GetAllSuspendedDistributorStatuses())

func ParseSuspendedDistributorStatus(status string) (string, error) {
	for _, s := range suspendedDistributorStatuses {
		if s == status {
			return s, nil
		}
	}

	return "", fmt.Errorf("'%s', %w", status, ErrorSuspendedDistributorStatusNotSupported)
}

func GetAllSuspendedDistributorStatuses() []string {
	return suspendedDistributorStatuses
}
