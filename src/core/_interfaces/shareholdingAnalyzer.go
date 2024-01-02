package _interfaces

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type ShareholdingAnalyzer interface {
	Analyze() (*entity2.RuleResultV2, *entity2.OwnershipStructure, error)
}
