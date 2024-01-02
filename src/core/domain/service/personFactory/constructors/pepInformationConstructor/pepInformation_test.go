package pepInformationConstructor

import (
	"testing"

	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Assemble_Should_Get_Pep_Information(t *testing.T) {
	restrictiveListsAdapter := &mocks.RestrictiveListsAdapter{}
	constructor := New(restrictiveListsAdapter)

	person := entity.Person{
		Name:           "name",
		EntityID:       uuid.New(),
		DocumentNumber: "123",
		PartnerID:      "SomePartnerID",
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ValidationSteps: []entity.ValidationStep{
				{
					RulesConfig: &entity.RuleSetConfig{
						PepParams: &entity.PepParams{},
					},
				},
			},
		},
	}
	personWrapper := &entity.PersonWrapper{
		Person: person,
	}
	pepInformation := &entity.PepInformation{
		DocumentNumber: person.DocumentNumber,
	}

	restrictiveListsAdapter.On("OccurrenceInPepList", person.DocumentNumber).Return(pepInformation, nil)

	err := constructor.Assemble(personWrapper)

	expected := pepInformation

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.PEPInformation)
	mock.AssertExpectationsForObjects(t, restrictiveListsAdapter)
}

func Test_Assemble_Should_Not_Get_Pep_Information(t *testing.T) {
	restrictiveListsAdapter := &mocks.RestrictiveListsAdapter{}
	constructor := New(restrictiveListsAdapter)

	person := entity.Person{
		Name:           "name",
		EntityID:       uuid.New(),
		DocumentNumber: "123",
		PartnerID:      "SomePartnerID",
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ValidationSteps: []entity.ValidationStep{
				{
					RulesConfig: &entity.RuleSetConfig{
						PepParams: &entity.PepParams{},
					},
				},
			},
		},
	}
	personWrapper := &entity.PersonWrapper{
		Person: person,
	}
	var pepInformation *entity.PepInformation = nil

	restrictiveListsAdapter.On("OccurrenceInPepList", person.DocumentNumber).Return(pepInformation, nil)

	err := constructor.Assemble(personWrapper)

	assert.Nil(t, err)
	assert.Equal(t, pepInformation, personWrapper.Person.PEPInformation)
	mock.AssertExpectationsForObjects(t, restrictiveListsAdapter)
}

func Test_Assemble_Should_Not_Validate_Pep(t *testing.T) {
	restrictiveListsAdapter := &mocks.RestrictiveListsAdapter{}
	constructor := New(restrictiveListsAdapter)

	person := entity.Person{
		Name:           "name",
		EntityID:       uuid.New(),
		DocumentNumber: "123",
		PartnerID:      "SomePartnerID",
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ValidationSteps: []entity.ValidationStep{
				{
					RulesConfig: &entity.RuleSetConfig{},
				},
			},
		},
	}
	personWrapper := &entity.PersonWrapper{
		Person: person,
	}
	var pepInformation *entity.PepInformation = nil

	err := constructor.Assemble(personWrapper)

	assert.Nil(t, err)
	assert.Equal(t, pepInformation, personWrapper.Person.PEPInformation)
	mock.AssertExpectationsForObjects(t, restrictiveListsAdapter)
}
