package legalRepresentativesConstructor

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"testing"

	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Assemble_When_LR_Rule_Enabled_Should_Get_LRs(t *testing.T) {

	profileID := uuid.New()

	profile := entity.Profile{
		ProfileID: &profileID,
		Person: entity.Person{
			PersonType: values.PersonTypeCompany,
			CadastralValidationConfig: &entity.CadastralValidationConfig{
				ValidationSteps: []entity.ValidationStep{
					{
						RulesConfig: &entity.RuleSetConfig{
							LegalRepresentativeParams: &entity.LegalRepresentativeParams{},
						},
					},
				},
			},
		},
	}

	profileWrapper := entity.ProfileWrapper{
		Profile: profile,
	}

	legalRepresentatives := []entity.LegalRepresentative{{LegalRepresentativeID: uuid.New()}}

	legalRepresentativeAdapter := &mocks.LegalRepresentativeAdapter{}
	legalRepresentativeAdapter.On("Search", *profile.ProfileID).Return(legalRepresentatives, nil)

	constructor := legalRepresentativesConstructor{legalRepresentativeAdapter: legalRepresentativeAdapter}

	err := constructor.Assemble(&profileWrapper)

	expected := legalRepresentatives

	assert.Nil(t, err)
	assert.Equal(t, &expected, &profileWrapper.Profile.LegalRepresentatives)
	mock.AssertExpectationsForObjects(t, legalRepresentativeAdapter)

}

func Test_Assemble_When_LR_Rule_Disabled_Should_Not_Get_LRs(t *testing.T) {

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

	legalRepresentativeAdapter := &mocks.LegalRepresentativeAdapter{}

	constructor := legalRepresentativesConstructor{legalRepresentativeAdapter: legalRepresentativeAdapter}

	err := constructor.Assemble(&profileWrapper)

	var expected []entity.LegalRepresentative = nil

	assert.Nil(t, err)
	assert.Equal(t, expected, profileWrapper.Profile.LegalRepresentatives)
	mock.AssertExpectationsForObjects(t, legalRepresentativeAdapter)

}
