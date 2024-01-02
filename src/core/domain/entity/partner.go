package entity

import (
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/pkg/errors"
)

type Partner struct {
	PartnerID      string
	DocumentNumber string
	Name           string
	Status         values2.PartnerStatus
	LogoImageUrl   string
	Config         *PartnerConfig
}

type PartnerConfig struct {
	CustomerSegregationType values2.SegregationType
	UseCallbackV2           *bool
}

func (ref *Partner) Validate() error {
	if ref == nil {
		return values2.NewErrorValidation("Partner is required")
	}

	if err := ref.Status.Validate(); err != nil {
		return errors.Wrap(err, "Invalid StatusCode")
	}

	if err := ref.Config.Validate(); err != nil {
		return errors.Wrap(err, "Invalid Config")
	}

	return nil
}

func (ref *PartnerConfig) Validate() error {
	if ref == nil {
		return values2.NewErrorValidation("Config is required")
	}

	if ref.CustomerSegregationType == "" {
		return values2.NewErrorValidation("CustomerSegregationType is required")
	}

	err := ref.CustomerSegregationType.Validate()
	if err != nil {
		return values2.NewErrorValidation("CustomerSegregationType is not valid. Valid values : (BY_PARTNER,BY_MERCHANT)")
	}

	if ref.UseCallbackV2 == nil {
		return values2.NewErrorValidation("UseCallbackV2 is required")
	}

	return nil
}

func (ref *PartnerConfig) IsCustomerSegregatedBy() values2.SegregationType {
	if ref == nil {
		return values2.SegregationTypeByPartner
	}
	return ref.CustomerSegregationType
}

func (ref *PartnerConfig) IsToUseCallbackV2() bool {
	if ref == nil || ref.UseCallbackV2 == nil {
		return false
	}
	return *ref.UseCallbackV2
}
