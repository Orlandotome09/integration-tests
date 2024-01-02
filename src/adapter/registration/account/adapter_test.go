package account

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/account/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/account/http/contracts"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"os"
	"reflect"
	"testing"
)

var (
	accountclient *accountClient.MockAccountClient
	adapter       interfaces.AccountAdapter
)

func TestMain(m *testing.M) {
	accountclient = &accountClient.MockAccountClient{}
	adapter = NewAccountAdapter(accountclient)
	os.Exit(m.Run())
}

func TestGetByID(t *testing.T) {
	accountID := "111"
	resp := &contracts.AccountResponse{ProfileID: uuid.MustParse("1d2d0cc3-f173-4a78-86c8-5b9de2fbf8d1")}

	accountclient.On("Get", accountID).Return(resp, nil)

	expected := &entity.Account{ProfileID: &resp.ProfileID}
	received, err := adapter.GetByID(accountID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
}

func TestFindByProfileID(t *testing.T) {
	profileID := uuid.New().String()

	accountResponse := []contracts.AccountResponse{
		{AccountNumber: "x1", AccountDigit: "x2", AgencyNumber: "x3", AgencyDigit: "x4", BankCode: "x5"},
		{AccountNumber: "y1", AccountDigit: "y2", AgencyNumber: "y3", AgencyDigit: "y4", BankCode: "y5"},
	}

	accountclient.On("Find", profileID).Return(accountResponse, nil)

	expected := []entity.Account{
		{
			AccountNumber: accountResponse[0].AccountNumber,
			AccountDigit:  accountResponse[0].AccountDigit,
			AgencyNumber:  accountResponse[0].AgencyNumber,
			AgencyDigit:   accountResponse[0].AgencyDigit,
			BankCode:      accountResponse[0].BankCode,
		},
		{
			AccountNumber: accountResponse[1].AccountNumber,
			AccountDigit:  accountResponse[1].AccountDigit,
			AgencyNumber:  accountResponse[1].AgencyNumber,
			AgencyDigit:   accountResponse[1].AgencyDigit,
			BankCode:      accountResponse[1].BankCode,
		},
	}

	received, err := adapter.FindByProfileID(profileID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}
	mock.AssertExpectationsForObjects(t, accountclient)
}

func TestCreateInternal(t *testing.T) {
	entityId := "111"
	bankCode := "133"
	request := &contracts.CreateAccountRequest{ProfileID: entityId, BankCode: bankCode}
	resp := &contracts.AccountResponse{ProfileID: uuid.MustParse("1d2d0cc3-f173-4a78-86c8-5b9de2fbf8d1")}

	accountclient.On("CreateInternal", request).Return(resp, nil)

	expected := &entity.Account{ProfileID: &resp.ProfileID}
	received, err := adapter.CreateInternal(entityId, bankCode)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected err nil")
	}

}
