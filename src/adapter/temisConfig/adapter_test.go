package temisConfig

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/temisConfig/http"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCadastralValidationConfig(t *testing.T) {

	temisConfigHttpClient := &http.MockTemisConfigHttpClient{}
	adapter := NewTemisConfigAdapter(temisConfigHttpClient)

	ctx := context.TODO()
	personType := "INDIVIDUAL"
	roleType := "CUSTOMER"
	offerType := "OFFER_TEST"
	partnerID := uuid.NewString()
	productConfig := &http.ProductConfig{
		CreateBexsAccount:           true,
		EnrichProfileWithBureauData: true,
		TreeIntegration:             true,
		LimitIntegration:            true,
	}

	response := http.CadastralValidationConfigResponse{
		{
			PersonType:    personType,
			OfferType:     offerType,
			RoleType:      roleType,
			PartnerID:     partnerID,
			ProductConfig: productConfig,
		},
	}

	temisConfigHttpClient.On("GetCadastralValidationConfig", ctx, personType, roleType, offerType).Return(response, nil)

	expected := &entity.CadastralValidationConfig{
		OfferType:       offerType,
		RoleType:        roleType,
		PersonType:      personType,
		PartnerID:       partnerID,
		ValidationSteps: entity.ValidationSteps{},
		ProductConfig: &entity.ProductConfig{
			CreateBexsAccount:           productConfig.CreateBexsAccount,
			EnrichProfileWithBureauData: productConfig.EnrichProfileWithBureauData,
			TreeIntegration:             productConfig.TreeIntegration,
			LimitIntegration:            productConfig.LimitIntegration,
		},
	}

	received, err := adapter.GetCadastralValidationConfig(ctx, personType, roleType, offerType, partnerID)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)

}

func TestFindConfig(t *testing.T) {
	t.Run("should find by complete key", func(t *testing.T) {
		personType := "INDIVIDUAL"
		roleType := "CUSTOMER"
		offerType := "OFFER_TEST"
		partnerID := uuid.NewString()

		completeKey := personType + roleType + offerType + partnerID
		baseKey := personType + roleType + offerType

		config1 := entity.CadastralValidationConfig{}

		config := entity.CadastralValidationConfigMap{
			completeKey: config1,
		}

		expected := &config1
		received := findConfig(completeKey, baseKey, config)

		assert.Equal(t, expected, received)
	})

	t.Run("should find by base key when not found by complete key", func(t *testing.T) {
		personType := "INDIVIDUAL"
		roleType := "CUSTOMER"
		offerType := "OFFER_TEST"
		partnerID := uuid.NewString()

		completeKey := personType + roleType + offerType + partnerID
		baseKey := personType + roleType + offerType

		config1 := entity.CadastralValidationConfig{}
		config2 := entity.CadastralValidationConfig{}

		config := entity.CadastralValidationConfigMap{
			"xxxxx": config1,
			baseKey: config2,
		}

		expected := &config2
		received := findConfig(completeKey, baseKey, config)

		assert.Equal(t, expected, received)
	})

	t.Run("should not find config", func(t *testing.T) {
		personType := "INDIVIDUAL"
		roleType := "CUSTOMER"
		offerType := "OFFER_TEST"
		partnerID := uuid.NewString()

		completeKey := personType + roleType + offerType + partnerID
		baseKey := personType + roleType + offerType

		config1 := entity.CadastralValidationConfig{}
		config2 := entity.CadastralValidationConfig{}

		config := entity.CadastralValidationConfigMap{
			"xxxxx": config1,
			"yyyy":  config2,
		}

		var expected *entity.CadastralValidationConfig
		received := findConfig(completeKey, baseKey, config)

		assert.Equal(t, expected, received)
	})
}
