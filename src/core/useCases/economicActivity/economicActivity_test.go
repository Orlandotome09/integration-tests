package economicActivity

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	validate             *validator.Validate
	economicActivityRepo *mocks.EconomicActivityRepository
	service              interfaces.EconomicActivityService
)

func TestMain(m *testing.M) {
	validate = validator.New()
	economicActivityRepo = &mocks.EconomicActivityRepository{}
	service = NewEconomicActivityService(validate, economicActivityRepo)
	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	code := "XXXX"
	economicalActivity := &entity.EconomicActivity{
		CodeID:      code,
		Description: "Some desc",
		RiskValue:   false,
	}

	economicActivityRepo.On("Get", code).Return(economicalActivity, true, nil)

	expected := economicalActivity
	received, exists, err := service.Get(code)

	assert.Nil(t, err)
	assert.Equal(t, true, exists)
	assert.Equal(t, expected, received)
}
