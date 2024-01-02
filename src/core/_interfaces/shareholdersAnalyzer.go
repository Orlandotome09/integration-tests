package _interfaces

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type ShareholdersAnalyzer interface {
	Analyze(ownershipStructure entity2.OwnershipStructure) (*entity2.RuleResultV2, error)
}
