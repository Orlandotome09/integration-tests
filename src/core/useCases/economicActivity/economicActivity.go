package economicActivity

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/go-playground/validator/v10"
)

type economicActivityService struct {
	validate   *validator.Validate
	repository interfaces.EconomicActivityRepository
}

func NewEconomicActivityService(validate *validator.Validate, repository interfaces.EconomicActivityRepository) interfaces.EconomicActivityService {
	return &economicActivityService{
		validate:   validate,
		repository: repository,
	}
}

func (ref *economicActivityService) Get(code string) (record *entity.EconomicActivity, exists bool, err error) {
	return ref.repository.Get(code)
}
