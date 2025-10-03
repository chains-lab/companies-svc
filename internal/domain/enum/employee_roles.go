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

func CheckEmployeeRole(role string) error {
	for _, r := range employeeRoles {
		if r == role {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", role, ErrorEmployeeRoleNotSupported)
}

func GetAllEmployeeRoles() []string {
	return employeeRoles
}

var AllEmployeeRoles = map[string]uint8{
	EmployeeRoleOwner:     math.MaxUint8,
	EmployeeRoleAdmin:     2,
	EmployeeRoleModerator: 1,
}

// CompareEmployeeRoles compares two employee roles and returns:
// 1 if role1 is greater than role2,
// 0 if they are equal.
// -1 if role1 is less than role2,
func CompareEmployeeRoles(role1, role2 string) (int, error) {
	err := CheckEmployeeRole(role1)
	if err != nil {
		return -1, fmt.Errorf("parsing role1: %w", err)
	}

	err = CheckEmployeeRole(role2)
	if err != nil {
		return -1, fmt.Errorf("parsing role2: %w", err)
	}

	if AllEmployeeRoles[role1] > AllEmployeeRoles[role2] {
		return 1, nil
	} else if AllEmployeeRoles[role1] < AllEmployeeRoles[role2] {
		return -1, nil
	}

	return 0, nil
}
