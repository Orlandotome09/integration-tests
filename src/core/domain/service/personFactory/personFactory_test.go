package personFactory

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_Person_Factory_Should_Build_Success(t *testing.T) {

	overrideService := &mocks.OverrideService{}
	catalogConstructor := &mocks.PersonConstructor{}
	personRulesConstructor := &mocks.PersonConstructor{}
	bureauConstructor := &mocks.PersonConstructor{}
	watchlistConstructor := &mocks.PersonConstructor{}
	blacklistConstructor := &mocks.PersonConstructor{}
	addressConstructor := &mocks.PersonConstructor{}
	documentsConstructor := &mocks.PersonConstructor{}

	personFactoryInstance := New(overrideService,
		catalogConstructor,
		personRulesConstructor,
		bureauConstructor,
		[]interfaces.PersonConstructor{
			watchlistConstructor,
			blacklistConstructor,
			addressConstructor,
			documentsConstructor,
		},
	)

	personID := uuid.New()
	person := entity.Person{
		EntityID:       personID,
		DocumentNumber: "123",
	}

	personWrapper := entity.PersonWrapper{
		Person: person,
	}

	timeDelay := 2 * time.Second
	catalogConstructor.On("Assemble", &personWrapper).Return(nil)
	personRulesConstructor.On("Assemble", &personWrapper).Return(nil)
	bureauConstructor.On("Assemble", &personWrapper).Return(nil)
	watchlistConstructor.On("Assemble", &personWrapper).After(timeDelay).Return(nil)
	blacklistConstructor.On("Assemble", &personWrapper).After(timeDelay).Return(nil)
	addressConstructor.On("Assemble", &personWrapper).After(timeDelay).Return(nil)
	documentsConstructor.On("Assemble", &personWrapper).After(timeDelay).Return(nil)
	overrideService.On("FindByEntityID", personWrapper.Person.EntityID).Return(nil, nil)

	timer := time.Now()
	personBuilt, err := personFactoryInstance.Build(person)

	expected := person
	maxTime := 3 * time.Second

	assert.Nil(t, err)
	assert.Equal(t, &expected, personBuilt)
	assert.GreaterOrEqual(t, time.Since(timer), timeDelay)
	assert.LessOrEqual(t, time.Since(timer), maxTime)
	mock.AssertExpectationsForObjects(t, catalogConstructor, personRulesConstructor, bureauConstructor,
		watchlistConstructor, blacklistConstructor, addressConstructor, documentsConstructor, overrideService)
}

func Test_Person_Factory_Error_Some_Goroutine_Constructor(t *testing.T) {

	overrideService := &mocks.OverrideService{}
	catalogConstructor := &mocks.PersonConstructor{}
	personRulesConstructor := &mocks.PersonConstructor{}
	bureauConstructor := &mocks.PersonConstructor{}
	watchlistConstructor := &mocks.PersonConstructor{}
	blacklistConstructor := &mocks.PersonConstructor{}
	addressConstructor := &mocks.PersonConstructor{}
	documentsConstructor := &mocks.PersonConstructor{}

	personFactoryInstance := New(overrideService,
		catalogConstructor,
		personRulesConstructor,
		bureauConstructor,
		[]interfaces.PersonConstructor{
			watchlistConstructor,
			blacklistConstructor,
			addressConstructor,
			documentsConstructor,
		},
	)

	personID := uuid.New()
	person := entity.Person{
		EntityID:       personID,
		DocumentNumber: "123",
	}

	personWrapper := entity.PersonWrapper{
		Person: person,
	}

	errorConstructor := errors.New("error")

	catalogConstructor.On("Assemble", &personWrapper).Return(nil)
	bureauConstructor.On("Assemble", &personWrapper).Return(nil)
	watchlistConstructor.On("Assemble", &personWrapper).Return(nil)
	blacklistConstructor.On("Assemble", &personWrapper).Return(nil)
	addressConstructor.On("Assemble", &personWrapper).Return(errorConstructor)
	documentsConstructor.On("Assemble", &personWrapper).Return(nil)

	personBuilt, err := personFactoryInstance.Build(person)

	expected := errorConstructor.Error()
	assert.Equal(t, expected, err.Error())
	assert.Nil(t, personBuilt)
	mock.AssertExpectationsForObjects(t, catalogConstructor, personRulesConstructor, bureauConstructor,
		watchlistConstructor, blacklistConstructor, addressConstructor, documentsConstructor, overrideService)
}
