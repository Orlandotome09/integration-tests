package values

import (
	"fmt"
)

type SegregationType string

const (
	SegregationTypeByPartner  SegregationType = "BY_PARTNER"
	SegregationTypeByMerchant SegregationType = "BY_MERCHANT"
)

var validCustomerSegregationTypes = map[string]SegregationType{
	string(SegregationTypeByPartner):  SegregationTypeByPartner,
	string(SegregationTypeByMerchant): SegregationTypeByMerchant,
}

func (ref *SegregationType) Validate() error {
	value := string(*ref)
	if _, in := validCustomerSegregationTypes[value]; !in {
		return NewErrorValidation(fmt.Sprintf("%s is an invalid segregation type", value))
	}
	return nil
}
