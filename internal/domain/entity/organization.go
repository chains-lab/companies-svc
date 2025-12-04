package entity

import (
	"fmt"

	"github.com/biter777/countries"
	"github.com/chains-lab/organizations-svc/internal/domain/errx"
	"github.com/google/uuid"
)

type CountryCode [3]byte

func (cc CountryCode) String() string {
	return string(cc[:])
}

type Organization struct {
	ID          uuid.UUID   `json:"id"`
	CountryCode CountryCode `json:"country_code"`
	Type        string      `json:"type"`
	Status      string      `json:"status"`
	Name        string      `json:"name"`
	Icon        *string     `json:"icon,omitempty"`
}

func ValidateCreateOrganization(org Organization) error {
	switch {
	case org.ID == uuid.Nil:
		return errx.ErrorOrganizationInvalidData.Raise(fmt.Errorf("organization ID is nil"))
	case ValidateOrganizationType(org.Type) != nil:
		return ValidateOrganizationType(org.Type)
	case ValidateOrganizationStatus(org.Status) != nil:
		return ValidateOrganizationStatus(org.Status)
	case org.Name == "":
		return errx.ErrorOrganizationInvalidData.Raise(fmt.Errorf("cannot create organization with empty name"))
	}

	if len(org.CountryCode.String()) != 3 {
		return errx.ErrorOrganizationInvalidData.Raise(fmt.Errorf("country code must be 3 characters long"))
	}

	if countries.ByName(org.CountryCode.String()) == countries.Unknown {
		return errx.ErrorOrganizationInvalidData.Raise(fmt.Errorf("invalid country code: %s", org.CountryCode.String()))
	}

	return nil
}

func (o Organization) IsNil() bool {
	if o.ID == uuid.Nil {
		return true
	}

	return false
}

const (
	OrganizationTypeCompany = "company"
	OrganizationTypeCity    = "city"
)

var organizationTypes = []string{
	OrganizationTypeCompany,
	OrganizationTypeCity,
}

func GetAllOrganizationTypes() []string {
	return organizationTypes
}

var ErrorOrganizationTypeNotSupported = fmt.Errorf("organization type not supported must be one of: %v", GetAllOrganizationTypes())

func ValidateOrganizationType(orgType string) error {
	for _, t := range organizationTypes {
		if t == orgType {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", orgType, ErrorOrganizationTypeNotSupported)
}

const (
	OrganizationStatusActive   = "active"
	OrganizationStatusInactive = "inactive"
	OrganizationStatusBanned   = "banned"
)

var organizationStatuses = []string{
	OrganizationStatusActive,
	OrganizationStatusInactive,
	OrganizationStatusBanned,
}

func GetAllOrganizationStatuses() []string {
	return organizationStatuses
}

var ErrorOrganizationStatusNotSupported = fmt.Errorf("organization status not supported must be one of: %v", GetAllOrganizationStatuses())

func ValidateOrganizationStatus(status string) error {
	for _, s := range organizationStatuses {
		if s == status {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", status, ErrorOrganizationStatusNotSupported)
}
