package partnerTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/partner/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockPartnerTranslator struct {
	PartnerTranslator
	mock.Mock
}

func (ref *MockPartnerTranslator) Translate(response contracts.PartnerResponse) entity.Partner {
	args := ref.Called(response)
	return args.Get(0).(entity.Partner)
}
