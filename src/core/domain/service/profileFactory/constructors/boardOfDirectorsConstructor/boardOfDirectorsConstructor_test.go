package boardOfDirectorsConstructor

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Assemble_When_LR_Rule_Enabled_Should_Get_BoardOfDirectors(t *testing.T) {

	profileID := uuid.New()

	profile := entity.Profile{
		ProfileID: &profileID,
		Person: entity.Person{
			PersonType: values.PersonTypeCompany,
			CadastralValidationConfig: &entity.CadastralValidationConfig{
				ValidationSteps: []entity.ValidationStep{
					{
						RulesConfig: &entity.RuleSetConfig{
							BoardOfDirectorsParams: &entity.BoardOfDirectorsParams{},
						},
					},
				},
			},
		},
	}

	profileWrapper := entity.ProfileWrapper{
		Profile: profile,
	}

	directors := []entity.Director{{DirectorID: uuid.New()}}

	boardOfDirectorsAdapter := &mocks.BoardOfDirectorsAdapter{}
	boardOfDirectorsAdapter.On("Search", *profile.ProfileID).Return(directors, nil)

	constructor := boardOfDirectorsConstructor{boardOfDirectorsAdapter: boardOfDirectorsAdapter}

	err := constructor.Assemble(&profileWrapper)

	expected := directors

	assert.Nil(t, err)
	assert.Equal(t, &expected, &profileWrapper.Profile.BoardOfDirectors)
	mock.AssertExpectationsForObjects(t, boardOfDirectorsAdapter)

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

	boardOfDirectorsAdapter := &mocks.BoardOfDirectorsAdapter{}

	constructor := boardOfDirectorsConstructor{boardOfDirectorsAdapter: boardOfDirectorsAdapter}

	err := constructor.Assemble(&profileWrapper)

	var expected []entity.Director = nil

	assert.Nil(t, err)
	assert.Equal(t, expected, profileWrapper.Profile.BoardOfDirectors)
	mock.AssertExpectationsForObjects(t, boardOfDirectorsAdapter)

}
