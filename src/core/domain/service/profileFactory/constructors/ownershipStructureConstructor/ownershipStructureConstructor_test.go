package ownershipStructureConstructor

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Assemble_When_OwnershipStructure_Rule_Enabled_Should_Get_OwnershipStructure(t *testing.T) {

	profileID := uuid.New()

	profile := entity.Profile{
		ProfileID: &profileID,
		Person: entity.Person{
			DocumentNumber: "123",
			PartnerID:      "SomePartnerID",
			OfferType:      "SomeOffer",
			PersonType:     values.PersonTypeCompany,
			CadastralValidationConfig: &entity.CadastralValidationConfig{
				ValidationSteps: []entity.ValidationStep{
					{
						RulesConfig: &entity.RuleSetConfig{
							OwnershipStructureParams: &entity.OwnershipStructureParams{},
						},
					},
				},
			},
			EnrichedInformation: &entity.EnrichedInformation{},
		},
	}

	profileWrapper := entity.ProfileWrapper{
		Profile: profile,
	}

	ownershipStructure := entity.OwnershipStructure{FinalBeneficiariesCount: 1}
	ownershipStructureEnriched := entity.OwnershipStructure{FinalBeneficiariesCount: 2}

	ownershipStructureService := &mocks.OwnershipStructureService{}
	ownershipStructureService.On("GetManuallyFilled", profile.ProfileID.String()).Return(&ownershipStructure, nil)
	ownershipStructureService.On("GetEnriched", profile.DocumentNumber, profile.OfferType, profile.PartnerID).Return(&ownershipStructureEnriched, nil)

	constructor := ownershipStructureConstructor{ownershipStructureService: ownershipStructureService}

	err := constructor.Assemble(&profileWrapper)

	expected := profile
	expected.EnrichedInformation.EnrichedCompany.OwnershipStructure = &ownershipStructureEnriched
	expected.OwnershipStructure = &ownershipStructure

	assert.Nil(t, err)
	assert.Equal(t, &expected.EnrichedInformation.EnrichedCompany.OwnershipStructure, &profileWrapper.Profile.EnrichedInformation.EnrichedCompany.OwnershipStructure)
	assert.Equal(t, &expected.OwnershipStructure, &profileWrapper.Profile.OwnershipStructure)
	mock.AssertExpectationsForObjects(t, ownershipStructureService)

}

func Test_Assemble_When_OwnershipStructure_Rule_Disabled_Should_Not_Get_OwnershipStructure(t *testing.T) {

	profileID := uuid.New()

	profile := entity.Profile{
		ProfileID: &profileID,
		Person: entity.Person{
			PersonType: values.PersonTypeCompany,
			CadastralValidationConfig: &entity.CadastralValidationConfig{
				ValidationSteps: []entity.ValidationStep{
					{},
				},
			},
		},
	}

	profileWrapper := entity.ProfileWrapper{
		Profile: profile,
	}

	ownershipStructureService := &mocks.OwnershipStructureService{}

	constructor := ownershipStructureConstructor{ownershipStructureService: ownershipStructureService}

	err := constructor.Assemble(&profileWrapper)

	var expected *entity.OwnershipStructure = nil

	assert.Nil(t, err)
	assert.Equal(t, expected, profileWrapper.Profile.OwnershipStructure)
	mock.AssertExpectationsForObjects(t, ownershipStructureService)

}
