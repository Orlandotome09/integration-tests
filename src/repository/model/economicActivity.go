package model

import (
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type EconomicalActivity struct {
	CodeID      string `gorm:"primaryKey"`
	Description string `gorm:"not null"`
	RiskValue   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (ref *EconomicalActivity) ToDomain() *entity.EconomicActivity {
	return &entity.EconomicActivity{
		CodeID:      ref.CodeID,
		Description: ref.Description,
		RiskValue:   ref.RiskValue,
	}
}

func (ref *EconomicalActivity) FromDomain(domainRiskyActivity entity.EconomicActivity) {

	ref.CodeID = domainRiskyActivity.CodeID
	ref.Description = domainRiskyActivity.Description
	ref.RiskValue = domainRiskyActivity.RiskValue
}
