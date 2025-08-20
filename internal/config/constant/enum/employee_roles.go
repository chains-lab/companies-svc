package enum

import (
	"fmt"
	"math"
)

const (
	EmployeeRoleOwner     = "owner"
	EmployeeRoleAdmin     = "admin"
	EmployeeRoleModerator = "moderator"
)

var employeeRoles = []string{
	EmployeeRoleOwner,
	EmployeeRoleAdmin,
	EmployeeRoleModerator,
}

var ErrorEmployeeRoleNotSupported = fmt.Errorf("employee role not supported, must be one of: %v", GetAllEmployeeRoles())

func ParseEmployeeRole(role string) (string, error) {
	for _, r := range employeeRoles {
		if r == role {
			return r, nil
		}
	}

	return "", fmt.Errorf("'%s', %w", role, ErrorEmployeeRoleNotSupported)
}

func GetAllEmployeeRoles() []string {
	return employeeRoles
}

var AllEmployeeRoles = map[string]uint8{
	EmployeeRoleOwner:     math.MaxUint8,
	EmployeeRoleAdmin:     2,
	EmployeeRoleModerator: 1,
}

// ComparisonEmployeeRoles compares two employee roles and returns:
// 1 if role1 is greater than role2,
// 0 if they are equal.
// -1 if role1 is less than role2,
func ComparisonEmployeeRoles(role1, role2 string) (int, error) {
	r1, err := ParseEmployeeRole(role1)
	if err != nil {
		return 0, fmt.Errorf("parsing role1: %w", err)
	}

	r2, err := ParseEmployeeRole(role2)
	if err != nil {
		return 0, fmt.Errorf("parsing role2: %w", err)
	}

	if AllEmployeeRoles[r1] > AllEmployeeRoles[r2] {
		return 1, nil
	} else if AllEmployeeRoles[r1] < AllEmployeeRoles[r2] {
		return -1, nil
	}

	return 0, nil
}
