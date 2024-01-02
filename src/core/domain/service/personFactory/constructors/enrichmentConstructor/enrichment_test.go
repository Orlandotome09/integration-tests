package enrichmentConstructor

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Assemble_Should_Enrich(t *testing.T) {

	adapter := &mocks.EnricherAdapter{}

	constructor := enrichmentPersonConstructor{adapter}

	person := entity.Person{
		EntityID:       uuid.New(),
		DocumentNumber: "123",
		OfferType:      "SomeOffer",
		PartnerID:      "SomePartner",
		PersonType:     values.PersonTypeIndividual,
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ValidationSteps: []entity.ValidationStep{
				{
					RulesConfig: &entity.RuleSetConfig{
						CAFParams: &entity.CAFParams{},
					},
				},
			},
		},
	}

	personWrapper := &entity.PersonWrapper{
		Person: person,
	}

	enrichedPerson := entity.EnrichedInformation{Providers: []entity.Provider{{ProviderName: "CAF"}}}

	adapter.On("GetEnrichedPerson",
		person.DocumentNumber,
		person.ProfileID.String(),
		person.PersonType,
		person.OfferType,
		person.PartnerID,
		person.RoleType).Return(&enrichedPerson, nil)

	err := constructor.Assemble(personWrapper)

	assert.NoError(t, err)
	assert.Equal(t, personWrapper.Person.EnrichedInformation, &enrichedPerson)
	mock.AssertExpectationsForObjects(t, adapter)

}

func Test_Assemble_Should_Merge_And_Enrich(t *testing.T) {

	adapter := &mocks.EnricherAdapter{}

	constructor := enrichmentPersonConstructor{adapter}

	person := entity.Person{
		EntityID:       uuid.New(),
		DocumentNumber: "123",
		OfferType:      "SomeOffer",
		PartnerID:      "SomePartner",
		PersonType:     values.PersonTypeIndividual,
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ValidationSteps: []entity.ValidationStep{
				{
					RulesConfig: &entity.RuleSetConfig{
						CAFParams: &entity.CAFParams{},
					},
				},
			},
		},
		EnrichedInformation: &entity.EnrichedInformation{BureauStatus: "SOME_STATUS"},
	}

	personWrapper := &entity.PersonWrapper{
		Person: person,
	}

	enrichedPerson := entity.EnrichedInformation{Providers: []entity.Provider{{ProviderName: "CAF"}}}

	adapter.On("GetEnrichedPerson",
		person.DocumentNumber,
		person.ProfileID.String(),
		person.PersonType,
		person.OfferType,
		person.PartnerID,
		person.RoleType).Return(&enrichedPerson, nil)

	err := constructor.Assemble(personWrapper)

	expected := enrichedPerson
	expected.BureauStatus = "SOME_STATUS"

	assert.NoError(t, err)
	assert.Equal(t, personWrapper.Person.EnrichedInformation, &expected)
	mock.AssertExpectationsForObjects(t, adapter)

}

func Test_Assemble_Should_Not_Enrich(t *testing.T) {

	adapter := &mocks.EnricherAdapter{}

	constructor := enrichmentPersonConstructor{adapter}

	person := entity.Person{
		EntityID:       uuid.New(),
		DocumentNumber: "123",
		OfferType:      "SomeOffer",
		PartnerID:      "SomePartner",
		PersonType:     values.PersonTypeIndividual,
	}

	personWrapper := &entity.PersonWrapper{
		Person: person,
	}

	err := constructor.Assemble(personWrapper)

	assert.NoError(t, err)
	assert.Nil(t, personWrapper.Person.EnrichedInformation)
	mock.AssertExpectationsForObjects(t, adapter)

}
