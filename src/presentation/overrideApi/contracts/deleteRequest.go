package contracts

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
)

type DeleteRequest struct {
	EntityID   uuid.UUID  `json:"entity_id" binding:"required"`
	EntityType string     `json:"entity_type" binding:"required"`
	ParentID   *uuid.UUID `json:"parent_id"`
	RuleName   string     `json:"rule_name" binding:"required"`
	RuleSet    string     `json:"rule_set" binding:"required"`
}

func (ref DeleteRequest) ToDomain() entity.Override {
	var parentID = ref.ParentID
	if parentID == nil {
		parentID = &ref.EntityID
	}
	return entity.Override{
		EntityID:   ref.EntityID,
		EntityType: values2.EntityType(ref.EntityType),
		ParentID:   parentID,
		RuleName:   values2.RuleName(ref.RuleName),
		RuleSet:    values2.RuleSet(ref.RuleSet),
	}
}
