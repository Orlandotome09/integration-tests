package _interfaces

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"time"
)

type ComplianceAnalyzer interface {
	RunComplianceAnalysis(complianceValidator ComplianceValidator, entityID uuid.UUID, requestDate time.Time, executionTime time.Time,
		cacheValue bool) (*entity.State, error)
}
