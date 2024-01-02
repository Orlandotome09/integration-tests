package contracts

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type Scores []Score

type Score struct {
	FileID          uuid.UUID          `json:"file_id"`
	Total           float64            `json:"total"`
	ForDocumentType ScoredDocumentType `json:"for_document_type"`
	ForDocumentSide ScoredDocumentSide `json:"for_document_side"`
	ForFields       ScoredFields       `json:"for_fields"`
}

func (ref Scores) ToDomain() entity2.Scores {
	scores := []entity2.Score{}
	for _, score := range ref {
		newScore := entity2.Score{
			FileID:          score.FileID,
			Total:           score.Total,
			ForDocumentType: score.ForDocumentType.ToDomain(),
			ForDocumentSide: score.ForDocumentSide.ToDomain(),
			ForFields:       score.ForFields.ToDomain(),
		}

		scores = append(scores, newScore)
	}
	return scores

}

func (ref *Scores) FromDomain(d entity2.Scores) {
	scores := []Score{}

	for _, s := range d {
		score := Score{
			FileID:          s.FileID,
			Total:           s.Total,
			ForDocumentType: ScoredDocumentType{},
			ForDocumentSide: ScoredDocumentSide{},
			ForFields:       ScoredFields{},
		}

		(&score.ForDocumentType).FromDomain(s.ForDocumentType)
		(&score.ForDocumentSide).FromDomain(s.ForDocumentSide)
		(&score.ForFields).FromDomain(s.ForFields)

		scores = append(scores, score)
	}
	*ref = scores
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

func (ref ScoredDocumentType) ToDomain() entity2.ScoredDocumentType {
	return entity2.ScoredDocumentType{
		Given:     ref.Given,
		Extracted: ref.Extracted,
		Score:     ref.Score,
	}
}

func (ref *ScoredDocumentType) FromDomain(d entity2.ScoredDocumentType) {
	ref.Given = d.Given
	ref.Extracted = d.Extracted
	ref.Score = d.Score
}

func (ref ScoredDocumentSide) ToDomain() entity2.ScoredDocumentSide {
	return entity2.ScoredDocumentSide{
		Given:     ref.Given,
		Extracted: ref.Extracted,
		Score:     ref.Score,
	}
}

func (ref *ScoredDocumentSide) FromDomain(d entity2.ScoredDocumentSide) {
	ref.Given = d.Given
	ref.Extracted = d.Extracted
	ref.Score = d.Score
}
