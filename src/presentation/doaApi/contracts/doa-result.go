package contracts

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"time"
)

type DOAResult struct {
	ID          uuid.UUID `json:"id"`
	EntityID    uuid.UUID `json:"entity_id"`
	DocumentID  uuid.UUID `json:"document_id"`
	RequestDate time.Time `json:"request_date"`
	Status      DOAStatus `json:"validation_status"`
	Scores      Scores    `json:"scores"`
}

func (ref DOAResult) ToDomain() entity.DOAResult {
	domainResult := entity.DOAResult{
		ID:          ref.ID,
		EntityID:    ref.EntityID,
		DocumentID:  ref.DocumentID,
		RequestDate: ref.RequestDate,
		Status:      ref.Status.ToDomain(),
		Scores:      ref.Scores.ToDomain(),
	}

	return domainResult
}

func (ref *DOAResult) FromDomain(d entity.DOAResult) {
	ref.ID = d.ID
	ref.EntityID = d.EntityID
	ref.DocumentID = d.DocumentID
	ref.RequestDate = d.RequestDate
	ref.Status.FromDomain(d.Status)
	ref.Scores = Scores{}
	(&ref.Scores).FromDomain(d.Scores)
}

type DOAStatus string

const (
	DOAStatusValidating DOAStatus = "VALIDATING"
	DOAStatusDone       DOAStatus = "DONE"
	DOAStatusError      DOAStatus = "ERROR"
)

func (ref DOAStatus) ToDomain() values.DOAStatus {
	switch ref {
	case DOAStatusValidating:
		return values.DOAStatusValidating
	case DOAStatusDone:
		return values.DOAStatusDone
	case DOAStatusError:
		return values.DOAStatusError
	default:
		return ""
	}
}

func (ref DOAStatus) FromDomain(d values.DOAStatus) {
	switch d {
	case values.DOAStatusValidating:
		ref = DOAStatusValidating
		return

	case values.DOAStatusDone:
		ref = DOAStatusDone
		return

	case values.DOAStatusError:
		ref = DOAStatusError
		return
	default:
		return
	}
}
