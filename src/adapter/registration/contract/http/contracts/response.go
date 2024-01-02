package contracts

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
)

type GetContractResponse struct {
	ContractID           *uuid.UUID `json:"contract_id"`
	EstimatedTotalAmount float64    `json:"estimated_total_amount" binding:"required,numeric"`
	DueDate              string     `json:"due_date" binding:"required"`
	ProfileID            string     `json:"profile_id" binding:"required"`
	DocumentID           string     `json:"document_id" binding:"required"`
}

func (ref *GetContractResponse) ToDomain() (*entity.Contract, error) {

	profileID, err := uuid.Parse(ref.ProfileID)
	if err != nil {
		return nil, err
	}

	var invoiceID *uuid.UUID

	if ref.DocumentID != "" {
		documentID, err := uuid.Parse(ref.DocumentID)
		if err != nil {
			return nil, err
		}
		invoiceID = &documentID
	} else {
		invoiceID = nil
	}

	return &entity.Contract{
		ContractID:           ref.ContractID,
		EstimatedTotalAmount: ref.EstimatedTotalAmount,
		DueTime:              ref.DueDate,
		Installments:         0,
		CorrelationID:        "",
		ProfileID:            &profileID,
		DocumentID:           invoiceID,
	}, nil

}
