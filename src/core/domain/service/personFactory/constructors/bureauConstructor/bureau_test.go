package bureauConstructor

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_Assemble_When_Enrich_Is_Enabled_Should_Populate_Bureau(t *testing.T) {
	bureauService := &mocks.BureauService{}
	constructor := bureauPersonConstructor{bureauService: bureauService}

	person := entity2.Person{
		EntityID:       uuid.New(),
		DocumentNumber: "123",
		OfferType:      "SomeOffer",
		PartnerID:      "SomePartner",
		PersonType:     values.PersonTypeIndividual,
		CadastralValidationConfig: &entity2.CadastralValidationConfig{
			ProductConfig: &entity2.ProductConfig{
				EnrichProfileWithBureauData: true,
			},
		},
	}
	personWrapper := &entity2.PersonWrapper{
		Person: person,
	}
	bureauStatus := &entity2.EnrichedInformation{
		BureauStatus: "X",
		EnrichedIndividual: entity2.EnrichedIndividual{
			Name:      "ANA PAULA DA SILVA DE OLIVEIRA",
			BirthDate: "30/12/1980",
		},
	}

	bureauService.On("GetBureauStatus", person).Return(bureauStatus, nil)

	err := constructor.Assemble(personWrapper)

	expected := bureauStatus

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.EnrichedInformation)
	mock.AssertExpectationsForObjects(t, bureauService)
}

func Test_Assemble_When_High_Risk_Is_Enabled_Should_Populate_Bureau(t *testing.T) {
	bureauService := &mocks.BureauService{}
	constructor := bureauPersonConstructor{bureauService: bureauService}

	person := entity2.Person{
		EntityID:       uuid.New(),
		DocumentNumber: "123",
		OfferType:      "SomeOffer",
		PartnerID:      "SomePartner",
		PersonType:     values.PersonTypeIndividual,
		CadastralValidationConfig: &entity2.CadastralValidationConfig{
			ValidationSteps: []entity2.ValidationStep{
				{
					RulesConfig: &entity2.RuleSetConfig{
						ActivityRiskParams: &entity2.ActivityRiskParams{},
					},
				},
			},
		},
	}
	personWrapper := &entity2.PersonWrapper{
		Person: person,
	}
	bureauStatus := &entity2.EnrichedInformation{
		BureauStatus: "X",
	}

	bureauService.On("GetBureauStatus", person).Return(bureauStatus, nil)

	err := constructor.Assemble(personWrapper)

	expected := bureauStatus

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.EnrichedInformation)
	mock.AssertExpectationsForObjects(t, bureauService)
}

func Test_Assemble_When_Bureau_Rule_Is_Enabled_Should_Populate_Bureau(t *testing.T) {
	bureauService := &mocks.BureauService{}
	constructor := bureauPersonConstructor{bureauService: bureauService}

	person := entity2.Person{
		EntityID:       uuid.New(),
		DocumentNumber: "123",
		OfferType:      "SomeOffer",
		PartnerID:      "SomePartner",
		PersonType:     values.PersonTypeIndividual,
		CadastralValidationConfig: &entity2.CadastralValidationConfig{
			ValidationSteps: []entity2.ValidationStep{
				{
					RulesConfig: &entity2.RuleSetConfig{
						BureauParams: &entity2.BureauParams{},
					},
				},
			},
		},
	}
	personWrapper := &entity2.PersonWrapper{
		Person: person,
	}
	bureauStatus := &entity2.EnrichedInformation{
		BureauStatus: "X",
	}

	bureauService.On("GetBureauStatus", person).Return(bureauStatus, nil)

	err := constructor.Assemble(personWrapper)

	expected := bureauStatus

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.EnrichedInformation)
	mock.AssertExpectationsForObjects(t, bureauService)
}

func Test_Assemble_Should_Not_Populate_Bureau(t *testing.T) {
	bureauService := &mocks.BureauService{}
	constructor := bureauPersonConstructor{bureauService: bureauService}

	person := entity2.Person{
		EntityID:       uuid.New(),
		DocumentNumber: "123",
		OfferType:      "SomeOffer",
		PartnerID:      "SomePartner",
		PersonType:     values.PersonTypeIndividual,
	}
	personWrapper := &entity2.PersonWrapper{
		Person: person,
	}

	err := constructor.Assemble(personWrapper)

	var expected *entity2.EnrichedInformation = nil

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.EnrichedInformation)
	mock.AssertExpectationsForObjects(t, bureauService)

}
