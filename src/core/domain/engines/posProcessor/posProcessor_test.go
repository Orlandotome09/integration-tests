package posProcessor

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"os"
	"reflect"
	"testing"
)

var (
	partnerAdapter       *mocks.PartnerAdapter
	treeAdapterPublisher *mocks.TreeAdapterPublisher
	limitPublisher       *mocks.LimitPublisher
	accountAdapter       *mocks.AccountAdapter
	addressAdapter       *mocks.AddressService
	posProcessorService  interfaces.PosProcessor
)

func TestMain(m *testing.M) {
	partnerAdapter = &mocks.PartnerAdapter{}
	treeAdapterPublisher = &mocks.TreeAdapterPublisher{}
	limitPublisher = &mocks.LimitPublisher{}
	accountAdapter = &mocks.AccountAdapter{}
	addressAdapter = &mocks.AddressService{}
	posProcessorService = NewPosProcessor(partnerAdapter, accountAdapter, addressAdapter, treeAdapterPublisher, limitPublisher)

	os.Exit(m.Run())
}

func TestSendToTreeAdapter_should_send_profile_customer(t *testing.T) {
	profileID := uuid.New()
	profile := &entity.Profile{
		ProfileID: &profileID,
		Person: entity.Person{
			PartnerID: "xxx",
			RoleType:  values.RoleTypeCustomer,
		},
	}
	partner := &entity.Partner{Name: "myPartner"}
	accountID := uuid.New()
	accounts := []entity.Account{{AccountID: &accountID}}

	addresses := []entity.Address{}
	profileUpdated := entity.Profile{
		ProfileID: &profileID,
		Person: entity.Person{
			PartnerID: partner.Name,
			RoleType:  values.RoleTypeCustomer,
		},
	}

	catalogConfig := &entity.ProductConfig{TreeIntegration: true}

	partnerAdapter.On("GetActive", profile.PartnerID).Return(partner, nil)
	accountAdapter.On("FindByProfileID", profile.ProfileID.String()).Return(accounts, nil)
	addressAdapter.On("Search", profile.ProfileID.String()).Return(addresses, nil)
	treeAdapterPublisher.On("Send", profileUpdated, accounts, addresses).Return(nil)

	var expected error = nil
	received := posProcessorService.SendToTreeAdapter(*profile, *catalogConfig)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	mock.AssertExpectationsForObjects(t, partnerAdapter, treeAdapterPublisher)
}

func TestSendToTreeAdapter_should_not_send_profile_not_customer(t *testing.T) {
	profile := &entity.Profile{
		Person: entity.Person{
			PartnerID: "xxx",
			RoleType:  values.RoleTypeMerchant,
		},
	}

	catalogConfig := &entity.ProductConfig{TreeIntegration: false}

	var expected error = nil
	received := posProcessorService.SendToTreeAdapter(*profile, *catalogConfig)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	partnerAdapter.AssertNotCalled(t, "GetActive")
	accountAdapter.AssertNotCalled(t, "FindByProfileID")
	treeAdapterPublisher.AssertNotCalled(t, "Send")
}

func TestSendToTreeAdapter_should_not_send_profile_for_disable_tree_integration(t *testing.T) {
	profile := &entity.Profile{
		Person: entity.Person{
			PartnerID: "xxx",
			RoleType:  values.RoleTypeCustomer,
		},
	}

	catalogConfig := &entity.ProductConfig{TreeIntegration: false}

	var expected error = nil
	received := posProcessorService.SendToTreeAdapter(*profile, *catalogConfig)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	partnerAdapter.AssertNotCalled(t, "GetActive")
	accountAdapter.AssertNotCalled(t, "FindByProfileID")
	treeAdapterPublisher.AssertNotCalled(t, "Send")
}

func TestSendToLimit_should_not_send_profile_for_disable_limit_integration(t *testing.T) {
	profile := &entity.Profile{
		Person: entity.Person{
			PartnerID: "xxx",
			RoleType:  values.RoleTypeCustomer,
		},
	}

	state := &entity.State{Result: values.ResultStatusApproved}

	catalogConfig := &entity.ProductConfig{LimitIntegration: false}

	var expected error = nil
	received := posProcessorService.SendToLimit(*profile, *state, *catalogConfig)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	limitPublisher.AssertNotCalled(t, "Send")
}

func TestSendToLimit_should_send_profile_customer(t *testing.T) {
	profile := &entity.Profile{
		Person: entity.Person{
			PartnerID: "xxx",
			RoleType:  values.RoleTypeCustomer,
		},
	}

	profileUpdated := &entity.Profile{
		Person: entity.Person{
			PartnerID: profile.PartnerID,
			RoleType:  values.RoleTypeCustomer,
		},
	}

	state := &entity.State{Result: values.ResultStatusApproved}

	catalogConfig := &entity.ProductConfig{LimitIntegration: true}

	limitPublisher.On("Send", *profileUpdated, *state).Return(nil)

	var expected error = nil
	received := posProcessorService.SendToLimit(*profile, *state, *catalogConfig)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	mock.AssertExpectationsForObjects(t, partnerAdapter, treeAdapterPublisher)
}

func TestCreateInternalAccount(t *testing.T) {
	entityID := uuid.New()
	account := &entity.Account{}

	accountAdapter.On("CreateInternal", entityID.String(), BexsInternalBankCode).Return(account, nil)

	var expected error = nil
	received := posProcessorService.CreateInternalAccount(entityID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	mock.AssertExpectationsForObjects(t, accountAdapter)
}
