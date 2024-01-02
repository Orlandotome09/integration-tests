package watchlistConstructor

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Assemble_ShouldGetWatchlist_For_Individual(t *testing.T) {
	watchlistService := &mocks.WatchlistAdapter{}
	constructor := watchlistConstructor{watchlistService: watchlistService}

	dateOfBirth := time.Now()
	person := entity.Person{
		EntityID:       uuid.New(),
		Name:           "JOAO SILVA",
		DocumentNumber: "123",
		PersonType:     values.PersonTypeIndividual,
		Individual:     &entity.Individual{DateOfBirth: &dateOfBirth, Nationality: "BR"},
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ValidationSteps: []entity.ValidationStep{
				{
					RulesConfig: &entity.RuleSetConfig{
						WatchListParams: &entity.WatchListParams{},
					},
				},
			},
		},
	}
	personWrapper := &entity.PersonWrapper{
		Person: person,
	}
	watchlist := &entity.Watchlist{}

	watchlistService.On("SearchIndividual",
		person.DocumentNumber,
		"JOAO",
		"SILVA",
		"JOAO SILVA",
		"BR",
		person.Individual.DateOfBirth.Year()).Return(watchlist, nil)

	err := constructor.Assemble(personWrapper)

	expected := watchlist

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.Watchlist)
	mock.AssertExpectationsForObjects(t, watchlistService)
	watchlistService.AssertNumberOfCalls(t, "SearchCompany", 0)
}

func Test_Assemble_ShouldNotGetWatchlist_For_Individual_When_NotValidIndividualWatchlist(t *testing.T) {
	watchlistService := &mocks.WatchlistAdapter{}
	constructor := watchlistConstructor{watchlistService: watchlistService}

	person := entity.Person{
		EntityID:       uuid.New(),
		DocumentNumber: "123",
		PersonType:     values.PersonTypeIndividual,
		Individual:     &entity.Individual{},
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ValidationSteps: []entity.ValidationStep{
				{
					RulesConfig: &entity.RuleSetConfig{
						WatchListParams: &entity.WatchListParams{},
					},
				},
			},
		},
	}
	personWrapper := entity.PersonWrapper{
		Person: person,
	}

	err := constructor.Assemble(&personWrapper)

	var expected *entity.Watchlist = nil

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.Watchlist)
	mock.AssertExpectationsForObjects(t, watchlistService)
	watchlistService.AssertNumberOfCalls(t, "SearchIndividual", 0)
	watchlistService.AssertNumberOfCalls(t, "SearchCompany", 0)
}

func Test_Assemble_ShouldGetWatchlist_For_Company(t *testing.T) {
	watchlistService := &mocks.WatchlistAdapter{}
	constructor := watchlistConstructor{watchlistService: watchlistService}

	person := entity.Person{
		EntityID:       uuid.New(),
		DocumentNumber: "123",
		PersonType:     values.PersonTypeCompany,
		Company:        &entity.Company{LegalName: "SomeName"},
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ValidationSteps: []entity.ValidationStep{
				{
					RulesConfig: &entity.RuleSetConfig{
						WatchListParams: &entity.WatchListParams{},
					},
				},
			},
		},
	}
	personWrapper := entity.PersonWrapper{
		Person: person,
	}
	watchlist := entity.Watchlist{}

	watchlistService.On("SearchCompany",
		person.Company.LegalName,
		"BR").Return(&watchlist, nil)

	err := constructor.Assemble(&personWrapper)

	expected := &watchlist

	assert.Nil(t, err)
	assert.Equal(t, &expected, &personWrapper.Person.Watchlist)
	mock.AssertExpectationsForObjects(t, watchlistService)
	watchlistService.AssertNumberOfCalls(t, "SearchIndividual", 0)
}

func Test_Assemble_ShouldNotGetWatchlist_For_Company_When_NotValidCompanyWatchlist(t *testing.T) {
	watchlistService := &mocks.WatchlistAdapter{}
	constructor := watchlistConstructor{watchlistService: watchlistService}

	person := entity.Person{
		EntityID:       uuid.New(),
		DocumentNumber: "123",
		PersonType:     values.PersonTypeCompany,
		Company:        &entity.Company{},
		CadastralValidationConfig: &entity.CadastralValidationConfig{
			ValidationSteps: []entity.ValidationStep{
				{
					RulesConfig: &entity.RuleSetConfig{
						WatchListParams: &entity.WatchListParams{},
					},
				},
			},
		},
	}
	personWrapper := entity.PersonWrapper{
		Person: person,
	}

	err := constructor.Assemble(&personWrapper)

	var expected *entity.Watchlist = nil

	assert.Nil(t, err)
	assert.Equal(t, expected, personWrapper.Person.Watchlist)
	watchlistService.AssertNumberOfCalls(t, "SearchIndividual", 0)
	watchlistService.AssertNumberOfCalls(t, "SearchCompany", 0)
}

func Test_ExtractName_ShouldGetFirstAndLastName_For_CompositeNames(t *testing.T) {
	names := [5]string{
		"JOAO",
		"JOAO SILVA",
		"JOAO DA SILVA",
		"JOAO DA SILVA OLIVEIRA",
		"JOAO PAULO DA SILVA OLIVEIRA",
	}

	expectedNames := []entity.Person{
		entity.Person{Individual: &entity.Individual{FirstName: "JOAO", LastName: ""}},
		entity.Person{Individual: &entity.Individual{FirstName: "JOAO", LastName: "SILVA"}},
		entity.Person{Individual: &entity.Individual{FirstName: "JOAO", LastName: "SILVA"}},
		entity.Person{Individual: &entity.Individual{FirstName: "JOAO", LastName: "OLIVEIRA"}},
		entity.Person{Individual: &entity.Individual{FirstName: "JOAO", LastName: "OLIVEIRA"}},
	}

	var result []entity.Person
	for _, fullName := range names {
		firstName, lastName := extractFirstLastName(fullName)
		result = append(result, entity.Person{Individual: &entity.Individual{FirstName: firstName, LastName: lastName}})
	}

	assert.Equal(t, expectedNames, result)

}

func Test_GetFirstLastName_ShouldGetFirstAndLastName_From_EnrichedInformation(t *testing.T) {
	person := entity.Person{
		Name: "JOAO DA SILVA",
		EnrichedInformation: &entity.EnrichedInformation{
			EnrichedIndividual: entity.EnrichedIndividual{
				Name: "JOAO DA SILVA OLIVEIRA",
			},
		},
	}
	personWrapper := &entity.PersonWrapper{
		Person: person,
	}

	expectedFirstName := "JOAO"
	expectedLastName := "OLIVEIRA"
	expectedFullName := "JOAO DA SILVA OLIVEIRA"
	firstName, lastName, fullName := getNames(personWrapper)

	assert.Equal(t, expectedFirstName, firstName)
	assert.Equal(t, expectedLastName, lastName)
	assert.Equal(t, expectedFullName, fullName)

}

func Test_GetFirstLastName_ShouldGetFirstAndLastName_From_Registration(t *testing.T) {
	person := entity.Person{
		Name: "JOAO DA SILVA",
		EnrichedInformation: &entity.EnrichedInformation{
			EnrichedIndividual: entity.EnrichedIndividual{
				Name: "",
			},
		},
	}
	personWrapper := &entity.PersonWrapper{
		Person: person,
	}

	expectedFirstName := "JOAO"
	expectedLastName := "SILVA"
	expectedFullName := "JOAO DA SILVA"
	firstName, lastName, fullName := getNames(personWrapper)

	assert.Equal(t, expectedFirstName, firstName)
	assert.Equal(t, expectedLastName, lastName)
	assert.Equal(t, expectedFullName, fullName)
}

func Test_GetFirstLastName_ShouldGetEmptyFirstAndLastName_From_EmptyInformations(t *testing.T) {
	person := entity.Person{
		Name: "",
		EnrichedInformation: &entity.EnrichedInformation{
			EnrichedIndividual: entity.EnrichedIndividual{
				Name: "",
			},
		},
	}
	personWrapper := &entity.PersonWrapper{
		Person: person,
	}

	expectedFirstName := ""
	expectedLastName := ""
	expectedFullName := ""
	firstName, lastName, fullName := getNames(personWrapper)

	assert.Equal(t, expectedFirstName, firstName)
	assert.Equal(t, expectedLastName, lastName)
	assert.Equal(t, expectedFullName, fullName)

}
