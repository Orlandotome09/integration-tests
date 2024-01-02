package values

import (
	"fmt"
)

type PersonType = string

const (
	PersonTypeIndividual PersonType = "INDIVIDUAL"
	PersonTypeCompany    PersonType = "COMPANY"
)

var validProfileTypes = map[string]PersonType{
	PersonTypeIndividual: PersonTypeIndividual,
	PersonTypeCompany:    PersonTypeCompany,
}

func ParseToPersonType(value string) (PersonType, error) {
	if _, exists := validProfileTypes[value]; !exists {
		return "", NewErrorValidation(fmt.Sprintf("%s is an invalid profile type", value))
	}
	return PersonType(value), nil
}
