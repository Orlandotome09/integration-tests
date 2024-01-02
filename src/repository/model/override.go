package model

import (
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
)

type Override struct {
	EntityID   uuid.UUID `gorm:"primaryKey;index:idx_overrides_entity_id;type:uuid"`
	EntityType string
	RuleSet    string `gorm:"primaryKey"`
	RuleName   string `gorm:"primaryKey"`
	Result     string
	Author     string
	Comments   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (ref Override) ToDomain() entity.Override {
	return entity.Override{
		EntityID:   ref.EntityID,
		EntityType: values2.EntityType(ref.EntityType),
		RuleSet:    values2.RuleSet(ref.RuleSet),
		RuleName:   values2.RuleName(ref.RuleName),
		Result:     values2.Result(ref.Result),
		Author:     ref.Author,
		Comments:   ref.Comments,
		CreatedAt:  ref.CreatedAt,
	}
}

func (Override) FromDomain(input entity.Override) Override {
	return Override{
		CreatedAt:  input.CreatedAt,
		EntityID:   input.EntityID,
		EntityType: input.EntityType.ToString(),
		RuleSet:    input.RuleSet.ToString(),
		RuleName:   input.RuleName.ToString(),
		Result:     input.Result.ToString(),
		Author:     input.Author,
		Comments:   input.Comments,
	}
}
