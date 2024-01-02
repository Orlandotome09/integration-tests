package values

import (
	"fmt"
)

type RoleType = string

const (
	RoleTypeCustomer            RoleType = "CUSTOMER"
	RoleTypeMerchant            RoleType = "MERCHANT"
	RoleTypeBusinessPartner     RoleType = "BUSINESS_PARTNER"
	RoleTypeCounterParty        RoleType = "COUNTERPARTY"
	RoleTypeLegalRepresentative RoleType = "LEGAL_REPRESENTATIVE"
	RoleTypeShareholder         RoleType = "SHAREHOLDER"
	RoleTypeDirector            RoleType = "DIRECTOR"
)

var validRoleTypes = map[string]RoleType{
	RoleTypeCustomer:            RoleTypeCustomer,
	RoleTypeMerchant:            RoleTypeMerchant,
	RoleTypeBusinessPartner:     RoleTypeBusinessPartner,
	RoleTypeLegalRepresentative: RoleTypeLegalRepresentative,
	RoleTypeShareholder:         RoleTypeShareholder,
	RoleTypeCounterParty:        RoleTypeCounterParty,
	RoleTypeDirector:            RoleTypeDirector,
}

func ParseToRoleType(value string) (RoleType, error) {
	if _, exists := validRoleTypes[value]; !exists {
		return "", NewErrorValidation(fmt.Sprintf("%s is an invalid role type", value))
	}

	return value, nil
}
