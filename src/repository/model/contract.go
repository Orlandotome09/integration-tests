package model

import (
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type Contract struct {
	ContractID           uuid.UUID `gorm:"primaryKey;type:uuid"`
	EstimatedTotalAmount float64
	DueTime              string
	Installments         int
	CorrelationID        string
	ProfileID            uuid.UUID  `gorm:"type:uuid"`
	DocumentID           *uuid.UUID `gorm:"type:uuid"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func NewContractFromDomain(contract entity.Contract) Contract {
	return Contract{
		ContractID:           *contract.ContractID,
		EstimatedTotalAmount: contract.EstimatedTotalAmount,
		DueTime:              contract.DueTime,
		Installments:         contract.Installments,
		CorrelationID:        contract.CorrelationID,
		ProfileID:            *contract.ProfileID,
		DocumentID:           contract.DocumentID,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}

func (ref Contract) ToDomain() entity.Contract {
	return entity.Contract{
		ContractID:           &ref.ContractID,
		EstimatedTotalAmount: ref.EstimatedTotalAmount,
		DueTime:              ref.DueTime,
		Installments:         ref.Installments,
		CorrelationID:        ref.CorrelationID,
		ProfileID:            &ref.ProfileID,
		DocumentID:           ref.DocumentID,
		CreatedAt:            ref.CreatedAt,
		UpdatedAt:            ref.UpdatedAt,
	}
}
