package entity

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"time"
)

type DOAResult struct {
	ID          uuid.UUID
	EntityID    uuid.UUID
	DocumentID  uuid.UUID
	FileIDs     []uuid.UUID
	EngineName  string
	RequestDate time.Time
	Status      values.DOAStatus
	Scores      Scores
}

type Scores []Score

type Score struct {
	FileID          uuid.UUID
	Total           float64
	ForDocumentType ScoredDocumentType
	ForDocumentSide ScoredDocumentSide
	ForFields       []ScoredField
}

type ScoredField struct {
	FieldName      string
	GivenValue     string
	ExtractedValue string
	Score          float64
}

type ScoredDocumentType struct {
	Given     string
	Extracted string
	Score     float64
}
type ScoredDocumentSide struct {
	Given     string
	Extracted string
	Score     float64
}
