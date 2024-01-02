package values

import (
	"fmt"
)

type SubEngineName = string

const (
	SubEngineNameLegalRepresentative = "LEGAL_REPRESENTATIVE"
	SubEngineNameShareholder         = "SHAREHOLDER"
)

var validSubEngineNames = map[string]SubEngineName{
	SubEngineNameLegalRepresentative: SubEngineNameLegalRepresentative,
	SubEngineNameShareholder:         SubEngineNameShareholder,
}

func ParseToSubEngineName(value string) (SubEngineName, error) {
	parsed, exists := validSubEngineNames[value]
	if !exists {
		return "", NewErrorValidation(fmt.Sprintf("%s is an invalid sub-engine name", value))
	}
	return parsed, nil
}
