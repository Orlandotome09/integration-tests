package contracts

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type Phone struct {
	Type        string `json:"type"`
	AreaCode    string `json:"area_code"`
	CountryCode string `json:"country_code"`
	Number      string `json:"number"`
}

type Phones []Phone

func (ref *Phone) ToDomain(target *entity.Phone) {
	target.Type = ref.Type
	target.Type = ref.Type
	target.CountryCode = ref.CountryCode
	target.AreaCode = ref.AreaCode
	target.Number = ref.Number
}

func (ref *Phone) FromDomain(phone *entity.Phone) {
	ref.Type = phone.Type
	ref.CountryCode = phone.CountryCode
	ref.AreaCode = phone.AreaCode
	ref.Number = phone.Number
}

func (ref Phones) ToDomain() []entity.Phone {
	output := make([]entity.Phone, len(ref))

	for index := range ref {
		ref[index].ToDomain(&output[index])
	}

	return output
}

func NewPhonesFromDomain(dom []entity.Phone) Phones {
	output := make(Phones, len(dom))

	for index := range dom {
		output[index].FromDomain(&dom[index])
	}

	return output
}
