package ownershipStructure

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	adapterEnrichment   *mocks.OwnershipStructureAdapter
	adapterRegistration *mocks.OwnershipStructureAdapter
)

func TestGetEnriched(t *testing.T) {
	adapterEnrichment := &mocks.OwnershipStructureAdapter{}
	adapterRegistration := &mocks.OwnershipStructureAdapter{}
	service := New(adapterEnrichment, adapterRegistration)
	legalEntityID := "111"
	offerType := "xxx"
	partnerID := "222"
	ownershipStructure := &entity.OwnershipStructure{FinalBeneficiariesCount: 3, ShareholdingSum: 100}

	adapterEnrichment.On("Get", legalEntityID, offerType, partnerID).Return(ownershipStructure, nil)

	expected := ownershipStructure
	received, err := service.GetEnriched(legalEntityID, offerType, partnerID)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}

func TestGetManuallyFilled(t *testing.T) {
	adapterEnrichment := &mocks.OwnershipStructureAdapter{}
	adapterRegistration := &mocks.OwnershipStructureAdapter{}
	service := New(adapterEnrichment, adapterRegistration)
	profileID := "111"
	ownershipStructure := &entity.OwnershipStructure{FinalBeneficiariesCount: 4, ShareholdingSum: 80}

	adapterRegistration.On("Get", profileID, "", "").Return(ownershipStructure, nil)

	expected := ownershipStructure
	received, err := service.GetManuallyFilled(profileID)

	assert.Nil(t, err)
	assert.Equal(t, expected, received)
}
