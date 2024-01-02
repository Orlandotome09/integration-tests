package values

import (
	"fmt"
)

type DOAStatus string

const (
	DOAStatusValidating DOAStatus = "VALIDATING"
	DOAStatusDone       DOAStatus = "DONE"
	DOAStatusError      DOAStatus = "ERROR"
)

var validDOAStatus = map[string]DOAStatus{
	DOAStatusValidating.ToString(): DOAStatusValidating,
	DOAStatusDone.ToString():       DOAStatusDone,
	DOAStatusError.ToString():      DOAStatusError,
}

func (doaStatus DOAStatus) Validate() error {
	_, in := validDOAStatus[doaStatus.ToString()]
	if !in {
		return NewErrorValidation(fmt.Sprintf("%s is an invalid rule set name", doaStatus))
	}
	return nil
}

func (doaStatus DOAStatus) ToString() string {
	return string(doaStatus)
}
