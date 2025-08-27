package enum

import "fmt"

const (
	BlockStatusActive   = "active"    // Distributor is active
	BlockStatusCanceled = "cancelled" // Distributor is canceled
)

var BlockStatuses = []string{
	BlockStatusActive,
	BlockStatusCanceled,
}

var ErrorBlockStatusNotSupported = fmt.Errorf("block distributor status must be one of: %s", GetAllBlockStatuses())

func ParseBlockStatus(status string) error {
	for _, s := range BlockStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", status, ErrorBlockStatusNotSupported)
}

func GetAllBlockStatuses() []string {
	return BlockStatuses
}
