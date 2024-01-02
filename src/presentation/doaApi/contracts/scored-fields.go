package contracts

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
)

type ScoredFields []ScoredField

type ScoredField struct {
	FieldName      string  `json:"field_name"`
	GivenValue     string  `json:"given_value"`
	ExtractedValue string  `json:"extracted_value"`
	Score          float64 `json:"score"`
}

func (ref ScoredFields) ToDomain() []entity2.ScoredField {

	scoredFields := []entity2.ScoredField{}

	for _, field := range ref {
		scoredField := entity2.ScoredField{
			FieldName:      field.FieldName,
			GivenValue:     field.GivenValue,
			ExtractedValue: field.ExtractedValue,
			Score:          field.Score,
		}
		scoredFields = append(scoredFields, scoredField)
	}

	return scoredFields
}

func (ref *ScoredFields) FromDomain(d []entity2.ScoredField) {
	scoredFields := []ScoredField{}

	for _, s := range d {
		scoredField := ScoredField{
			FieldName:      s.FieldName,
			GivenValue:     s.GivenValue,
			ExtractedValue: s.ExtractedValue,
			Score:          s.Score,
		}
		scoredFields = append(scoredFields, scoredField)
	}
	*ref = scoredFields
}
