package addressConstructor

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Assemble_When_Tree_Integration_Enabled_Should_Populate_Address(t *testing.T) {
	addressService := &mocks.AddressService{}
	constructor := addressPersonConstructor{addressAdapter: addressService}

	person := entity.Person{
		EntityID: uuid.New(),
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ProductConfig: &entity.ProductConfig{
				TreeIntegration: true,
			},
		},
	}
	personWrapper := &entity.PersonWrapper{
		Person: person,
	}
	addresses := []entity.Address{
		{Street: "Some address"},
	}

	addressService.On("Search", person.EntityID.String()).Return(addresses, nil)

	err := constructor.Assemble(personWrapper)

	expected := addresses

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.Addresses)
	mock.AssertExpectationsForObjects(t, addressService)
}

func Test_Assemble_When_Address_Required_Should_Populate_Address(t *testing.T) {
	addressService := &mocks.AddressService{}
	constructor := addressPersonConstructor{addressAdapter: addressService}

	person := entity.Person{
		EntityID: uuid.New(),
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ValidationSteps: []entity.ValidationStep{
				{
					RulesConfig: &entity.RuleSetConfig{
						IncompleteParams: &entity.IncompleteParams{AddressRequired: true},
					},
				},
			},
		},
	}
	personWrapper := &entity.PersonWrapper{
		Person: person,
	}
	addresses := []entity.Address{
		{Street: "Some address"},
	}

	addressService.On("Search", person.EntityID.String()).Return(addresses, nil)

	err := constructor.Assemble(personWrapper)

	expected := addresses

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.Addresses)
	mock.AssertExpectationsForObjects(t, addressService)
}

func Test_Assemble_Should_Not_Populate_Address(t *testing.T) {
	addressService := &mocks.AddressService{}
	constructor := addressPersonConstructor{addressAdapter: addressService}

	person := entity.Person{
		EntityID:  uuid.New(),
		Addresses: nil,
	}

	personWrapper := &entity.PersonWrapper{
		Person: person,
	}

	err := constructor.Assemble(personWrapper)

	var expected []entity.Address = nil

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.Addresses)
	mock.AssertExpectationsForObjects(t, addressService)
}
