package accountTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/treeAdapter/message"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"reflect"
	"testing"
)

func TestTranslate(t *testing.T) {
	translator := NewAccountTranslator()

	accounts := []entity.Account{
		{AccountNumber: "111", AccountDigit: "4", AgencyNumber: "222", AgencyDigit: "7", BankCode: "333"},
		{AccountNumber: "444", AccountDigit: "5", AgencyNumber: "555", AgencyDigit: "6", BankCode: "666"},
	}

	expected := message.Accounts{
		{Number: accounts[0].AccountNumber + accounts[0].AccountDigit, AgencyNumber: accounts[0].AgencyNumber + accounts[0].AgencyDigit, BankCode: accounts[0].BankCode},
		{Number: accounts[1].AccountNumber + accounts[1].AccountDigit, AgencyNumber: accounts[1].AgencyNumber + accounts[1].AgencyDigit, BankCode: accounts[1].BankCode},
	}

	received := translator.Translate(accounts)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}
}
