package entity

import (
	"fmt"
	"time"

	"github.com/chains-lab/organizations-svc/internal/domain/errx"
	"github.com/google/uuid"
)

const OwnerRoleName = "owner"        //TODO: default role for company type organizations
const TechLeadRoleName = "tech_lead" //TODO: default role for city type organizations

type Role struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	OrgType        string    `json:"organization_type"`
	Editable       bool      `json:"editable"`
	Unique         bool      `json:"unique"`
	Rank           uint64    `json:"rank"`
	Name           string    `json:"name"`
	Description    *string   `json:"description,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	CityPermissions         *RoleCityPermission         `json:"city_permissions,omitempty"`
	OrganizationPermissions *RoleOrganizationPermission `json:"organization_permissions,omitempty"`
}

func ValidateCreateRole(role Role) error {
	switch {
	case role.ID != uuid.Nil:
		return errx.ErrRoleInvalidData.Raise(fmt.Errorf("cannot create role with predefined ID"))
	case role.OrganizationID == uuid.Nil:
		return errx.ErrRoleInvalidData.Raise(fmt.Errorf("cannot create role with nil organization ID"))
	case ValidateOrganizationType(role.OrgType) != nil:
		return ValidateOrganizationType(role.OrgType)
	case role.Name == "":
		return errx.ErrRoleInvalidData.Raise(fmt.Errorf("cannot create role with empty name"))
	}

	return nil
}

func (role Role) IsNil() bool {
	if role.ID == uuid.Nil {
		return true
	}

	return false
}

type RoleCityPermission struct {
	RoleID uuid.UUID `json:"role_id"`

	UpdateCity       bool `json:"update_city"`
	UpdateCityStatus bool `json:"update_city_status"`
	UpdateCitySlug   bool `json:"update_city_slug"`

	AnswerToPetitions bool `json:"answer_to_petitions"`
	HidePetitions     bool `json:"hide_petitions"`
}

func (roleCityPermission RoleCityPermission) IsNil() bool {
	if roleCityPermission.RoleID == uuid.Nil {
		return true
	}

	return false
}

type RoleOrganizationPermission struct {
	RoleID uuid.UUID `json:"role_id"`

	UpdateOrganization       bool `json:"update_organization"`
	ChangeOrganizationStatus bool `json:"change_organization_status"`

	InviteMembers bool `json:"invite_members"`
	DeleteMembers bool `json:"delete_members"`
	UpdateMembers bool `json:"update_members"`

	ManageRoles bool `json:"manage_roles"`
}

func (roleOrganizationPermission RoleOrganizationPermission) IsNil() bool {
	if roleOrganizationPermission.RoleID == uuid.Nil {
		return true
	}

	return false
}
