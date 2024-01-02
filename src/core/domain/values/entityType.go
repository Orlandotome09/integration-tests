package values

import (
	"fmt"
)

type EntityType string

const (
	EntityTypeProfile             EntityType = "PROFILE"
	EntityTypeContract            EntityType = "CONTRACT"
	EntityTypeLegalRepresentative EntityType = "LEGAL_REPRESENTATIVE"
	EntityTypeShareholder         EntityType = "SHAREHOLDER"
	EntityTypeDirector            EntityType = "DIRECTOR"
	EntityTypeComplianceState     EntityType = "COMPLIANCE_STATE"
)

var validEntityTypes = map[string]EntityType{
	EntityTypeProfile.ToString():             EntityTypeProfile,
	EntityTypeContract.ToString():            EntityTypeContract,
	EntityTypeLegalRepresentative.ToString(): EntityTypeLegalRepresentative,
	EntityTypeShareholder.ToString():         EntityTypeShareholder,
	EntityTypeDirector.ToString():            EntityTypeDirector,
	EntityTypeComplianceState.ToString():     EntityTypeComplianceState,
}

func (entityType EntityType) Validate() error {
	_, in := validEntityTypes[entityType.ToString()]
	if !in {
		return NewErrorValidation(fmt.Sprintf("%s is an invalid entity type", entityType))
	}
	return nil
}

func (entityType EntityType) ToString() string {
	return string(entityType)
}
