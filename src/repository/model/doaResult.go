package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type DOAResult struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	EntityID    uuid.UUID `gorm:"type:uuid"`
	DocumentID  uuid.UUID `gorm:"type:uuid"`
	FileIDs     FileIDs   `gorm:"type:jsonb"`
	EngineName  string
	RequestDate time.Time
	SentPayload postgres.Jsonb `gorm:"type:jsonb"`
	Status      string
	Scores      Scores `gorm:"type:jsonb"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Scores []Score

type FileIDs []uuid.UUID

func (ref FileIDs) Value() (driver.Value, error) {
	result, err := json.Marshal(ref)
	return result, err
}

func (ref *FileIDs) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.Errorf("Failed to unmarshal JSONB value: %+v", value)
	}

	return json.Unmarshal(bytes, ref)
}

func (ref Scores) Value() (driver.Value, error) {
	result, err := json.Marshal(ref)
	return result, err
}

func (ref *Scores) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.Errorf("Failed to unmarshal JSONB value: %+v", value)
	}

	return json.Unmarshal(bytes, ref)
}

func (ref *DOAResult) ToDomain() entity.DOAResult {

	result := entity.DOAResult{
		ID:          ref.ID,
		EntityID:    ref.EntityID,
		DocumentID:  ref.DocumentID,
		FileIDs:     ref.FileIDs,
		EngineName:  ref.EngineName,
		RequestDate: ref.CreatedAt,
		Status:      values.DOAStatus(ref.Status),
		Scores:      ref.Scores.ToDomain(),
	}
	return result
}

func (ref *DOAResult) FromDomain(doaResult *entity.DOAResult) {

	payload := postgres.Jsonb{}

	ref.ID = doaResult.ID
	ref.EntityID = doaResult.EntityID
	ref.DocumentID = doaResult.DocumentID
	ref.FileIDs = doaResult.FileIDs
	ref.EngineName = doaResult.EngineName
	ref.SentPayload = payload
	ref.RequestDate = doaResult.RequestDate
	ref.Status = string(doaResult.Status)
	ref.Scores = Scores{}
	(&ref.Scores).FromDomain(doaResult.Scores)
}

func (ref *DOAResult) PrepareWithCurrent(current *DOAResult) {

	if current != nil {
		//Not mutable fields
		ref.ID = current.ID
		ref.DocumentID = current.DocumentID
		ref.FileIDs = current.FileIDs
		ref.EntityID = current.EntityID
		ref.EngineName = current.EngineName
		ref.RequestDate = current.RequestDate
		ref.SentPayload = current.SentPayload
	}

}

func (ref Scores) ToDomain() entity.Scores {

	scores := entity.Scores{}

	for _, score := range ref {
		domainScore := entity.Score{
			FileID: score.FileID,
			Total:  score.Total,
			ForDocumentType: entity.ScoredDocumentType{
				Given:     score.ForDocumentType.Given,
				Extracted: score.ForDocumentType.Extracted,
				Score:     score.ForDocumentType.Score,
			},
			ForDocumentSide: entity.ScoredDocumentSide{
				Given:     score.ForDocumentSide.Given,
				Extracted: score.ForDocumentSide.Extracted,
				Score:     score.ForDocumentSide.Score,
			},
			ForFields: []entity.ScoredField{},
		}

		scoredFields := []entity.ScoredField{}
		for _, field := range score.ForFields {
			scoredField := entity.ScoredField{
				FieldName:      field.FieldName,
				GivenValue:     field.GivenValue,
				ExtractedValue: field.ExtractedValue,
				Score:          field.Score,
			}
			scoredFields = append(scoredFields, scoredField)
		}

		domainScore.ForFields = scoredFields
		scores = append(scores, domainScore)
	}

	return scores
}

func (ref *Scores) FromDomain(scores entity.Scores) {
	modelScores := Scores{}

	for _, score := range scores {
		modelScore := Score{
			FileID: score.FileID,
			Total:  score.Total,
			ForDocumentType: ScoredDocumentType{
				Given:     score.ForDocumentType.Given,
				Extracted: score.ForDocumentType.Extracted,
				Score:     score.ForDocumentType.Score,
			},
			ForDocumentSide: ScoredDocumentSide{
				Given:     score.ForDocumentSide.Given,
				Extracted: score.ForDocumentSide.Extracted,
				Score:     score.ForDocumentSide.Score,
			},
		}

		scoredFields := []ScoredField{}
		for _, field := range score.ForFields {
			scoredField := ScoredField{
				FieldName:      field.FieldName,
				GivenValue:     field.GivenValue,
				ExtractedValue: field.ExtractedValue,
				Score:          field.Score,
			}
			scoredFields = append(scoredFields, scoredField)
		}

		modelScore.ForFields = scoredFields
		modelScores = append(modelScores, modelScore)
	}
	*ref = modelScores
}

type Score struct {
	FileID          uuid.UUID          `json:"file_id"`
	Total           float64            `json:"total"`
	ForDocumentType ScoredDocumentType `json:"for_document_type"`
	ForDocumentSide ScoredDocumentSide `json:"for_document_side"`
	ForFields       []ScoredField      `json:"for_fields"`
}

type ScoredDocumentType struct {
	Given     string  `json:"given"`
	Extracted string  `json:"extracted"`
	Score     float64 `json:"score"`
}
type ScoredDocumentSide struct {
	Given     string  `json:"given"`
	Extracted string  `json:"extracted"`
	Score     float64 `json:"score"`
}

type ScoredField struct {
	FieldName      string  `json:"field_name"`
	GivenValue     string  `json:"given_value"`
	ExtractedValue string  `json:"extracted_value"`
	Score          float64 `json:"score"`
}
