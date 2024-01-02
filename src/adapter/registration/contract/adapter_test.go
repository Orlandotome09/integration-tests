package contract

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contract/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/contract/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestGetByID(t *testing.T) {

	contractClientInstance := &contractClient.MockContractClient{}
	service := NewContractAdapter(contractClientInstance)

	contractID := uuid.New()
	profileID := uuid.New()
	invoiceID := uuid.New()
	resp := &contracts.GetContractResponse{ContractID: &contractID, DocumentID: invoiceID.String(), ProfileID: profileID.String(), EstimatedTotalAmount: 0}

	contractClientInstance.On("Get", &contractID).Return(resp, true, nil)

	expected := &entity.Contract{
		ContractID:           &contractID,
		EstimatedTotalAmount: 0,
		DueTime:              "",
		Installments:         0,
		CorrelationID:        "",
		ProfileID:            &profileID,
		DocumentID:           &invoiceID,
	}
	received, exists, err := service.Get(&contractID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if exists != true {
		t.Errorf("\nExpected exists true")
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}

}
