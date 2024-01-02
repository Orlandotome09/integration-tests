package temisConfig

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/temisConfig/http"
	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"context"
	"fmt"
	"github.com/pkg/errors"
)

type temisConfigAdapter struct {
	temisConfigHttpClient http.TemisConfigHttpClient
}

func NewTemisConfigAdapter(temisConfigHttpClient http.TemisConfigHttpClient) _interfaces.TemisConfigAdapter {
	return &temisConfigAdapter{
		temisConfigHttpClient: temisConfigHttpClient,
	}
}

func (ref *temisConfigAdapter) GetCadastralValidationConfig(ctx context.Context, personType, roleType, offerType, partnerID string) (*entity.CadastralValidationConfig, error) {
	res, err := ref.temisConfigHttpClient.GetCadastralValidationConfig(ctx, personType, roleType, offerType)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("[CadastralValidationConfigAdapter] %s", err.Error()))
	}
	if len(res) == 0 {
		return nil, nil
	}

	cadastralValidaionConfigMap := res.ToDomain()
	completeKey := personType + roleType + offerType + partnerID
	baseKey := personType + roleType + offerType

	config := findConfig(completeKey, baseKey, cadastralValidaionConfigMap)

	return config, nil
}

func findConfig(completeKey string, baseKey string, configMap entity.CadastralValidationConfigMap) *entity.CadastralValidationConfig {
	config, exists := configMap[completeKey]
	if !exists {
		config, exists = configMap[baseKey]
		if !exists {
			return nil
		}
	}

	return &config
}