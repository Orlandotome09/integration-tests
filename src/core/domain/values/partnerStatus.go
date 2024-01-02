package values

import (
	"fmt"
)

type PartnerStatus string

const (
	PartnerStatusActive PartnerStatus = "ACTIVE"
)

var validPartnerStatus = map[string]PartnerStatus{
	string(PartnerStatusActive): PartnerStatusActive,
}

func (ref *PartnerStatus) Validate() error {
	if ref == nil {
		return NewErrorValidation("Empty PartnerStatus")
	}
	value := string(*ref)
	if _, exists := validPartnerStatus[value]; !exists {
		return NewErrorValidation(fmt.Sprintf("%s is an invalid partner status", value))
	}
	return nil
}
