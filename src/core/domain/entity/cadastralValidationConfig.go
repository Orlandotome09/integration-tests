package entity

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type CadastralValidationConfig struct {
	OfferType       string            `json:"offer_type"`
	RoleType        values.RoleType   `json:"role_type"`
	PersonType      values.PersonType `json:"person_type"`
	PartnerID       string            `json:"partner_id,omitempty"`
	ValidationSteps ValidationSteps   `json:"validation_steps,omitempty"`
	ProductConfig   *ProductConfig    `json:"product_config,omitempty"`
}

type CadastralValidationConfigMap map[string]CadastralValidationConfig

type ProductConfig struct {
	CreateBexsAccount           bool `json:"create_bexs_account"`
	EnrichProfileWithBureauData bool `json:"enrich_profile_with_bureau_data"`
	TreeIntegration             bool `json:"tree_integration"`
	LimitIntegration            bool `json:"limit_integration"`
}

func (cadastralValidationConfig CadastralValidationConfig) HasProductConfig() bool {
	return cadastralValidationConfig.ProductConfig != nil
}

func (cadastralValidationConfig CadastralValidationConfig) HasInternalAccountCreation() bool {
	return cadastralValidationConfig.HasProductConfig() && cadastralValidationConfig.ProductConfig.CreateBexsAccount
}

func (cadastralValidationConfig CadastralValidationConfig) HasTreeIntegration() bool {
	return cadastralValidationConfig.HasProductConfig() && cadastralValidationConfig.ProductConfig.TreeIntegration
}

func (cadastralValidationConfig CadastralValidationConfig) HasLimitIntegration() bool {
	return cadastralValidationConfig.HasProductConfig() && cadastralValidationConfig.ProductConfig.LimitIntegration
}

func (cadastralValidationConfig CadastralValidationConfig) HasBureauEnrichment() bool {
	return cadastralValidationConfig.HasProductConfig() && cadastralValidationConfig.ProductConfig.EnrichProfileWithBureauData
}
