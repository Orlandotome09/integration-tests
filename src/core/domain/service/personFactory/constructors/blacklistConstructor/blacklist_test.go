package blacklistConstructor

import (
	"testing"

	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Assemble_Should_Get_Blacklist(t *testing.T) {
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
						BlackListParams: &entity.BlackListParams{},
					},
				},
			},
		},
	}
	personWrapper := &entity.PersonWrapper{
		Person: person,
	}
	blacklist := &entity.BlacklistStatus{
		Status: "X",
	}

	restrictiveListsAdapter.On("OccurrenceInBlackList", person.DocumentNumber, person.Name).Return(blacklist, nil)

	err := constructor.Assemble(personWrapper)

	expected := blacklist

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.BlacklistStatus)
	mock.AssertExpectationsForObjects(t, restrictiveListsAdapter)
}

func Test_Assemble_Should_Not_Get_Blacklist(t *testing.T) {
	restrictiveListsAdapter := &mocks.RestrictiveListsAdapter{}
	constructor := New(restrictiveListsAdapter)

	person := entity.Person{
		EntityID:       uuid.New(),
		DocumentNumber: "123",
		PartnerID:      "SomePartnerID",
	}
	personWrapper := &entity.PersonWrapper{
		Person: person,
	}

	err := constructor.Assemble(personWrapper)

	var expected *entity.BlacklistStatus = nil

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.BlacklistStatus)
	mock.AssertExpectationsForObjects(t, restrictiveListsAdapter)
}

func Test_Assemble_Should_Get_Empty_Blacklist(t *testing.T) {
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
						BlackListParams: &entity.BlackListParams{},
					},
				},
			},
		},
	}
	personWrapper := &entity.PersonWrapper{
		Person: person,
	}

	restrictiveListsAdapter.On("OccurrenceInBlackList", person.DocumentNumber, person.Name).Return(nil, nil)

	err := constructor.Assemble(personWrapper)

	var expected *entity.BlacklistStatus = nil

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.BlacklistStatus)
	mock.AssertExpectationsForObjects(t, restrictiveListsAdapter)
}
