package values

import (
	"fmt"
)

type EngineName = string

const (
	EngineNameProfile  = "PROFILE"
	EngineNameContract = "CONTRACT"
	EngineNamePerson   = "PERSON"
)

var validEngineNames = map[string]EngineName{
	EngineNameProfile:  EngineNameProfile,
	EngineNameContract: EngineNameContract,
	EngineNamePerson:   EngineNamePerson,
}

func ParseToEngineName(value string) (EngineName, error) {
	parsed, exists := validEngineNames[value]
	if !exists {
		return "", NewErrorValidation(fmt.Sprintf("%s is an invalid engine name", value))
	}
	return parsed, nil
}
