package boardOfDirectorsTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/boardOfDirectors/http/contracts"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/pkg/errors"
	"time"
)

type BoardOfDirectorsTranslator interface {
	Translate(response contracts.BoardOfDirectorsResponse) (*entity.Director, error)
}

type boardOfDirectorsTranslator struct{}

func New() BoardOfDirectorsTranslator {
	return &boardOfDirectorsTranslator{}
}

func (ref *boardOfDirectorsTranslator) Translate(response contracts.BoardOfDirectorsResponse) (*entity.Director, error) {

	var dateOfBirth *time.Time = nil
	if response.DateOfBirth != "" {
		result, err := time.Parse("2006-01-02", response.DateOfBirth)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		dateOfBirth = &result
	}

	return &entity.Director{
		DirectorID: response.DirectorID,
		Person: entity.Person{
			DocumentNumber: response.DocumentNumber,
			Name:           response.FullName,
			PersonType:     values.PersonTypeIndividual,
			ProfileID:      response.ProfileID,
			EntityID:       response.DirectorID,
			EntityType:     values.EntityTypeDirector,
			RoleType:       values.RoleTypeDirector,
			Individual: &entity.Individual{
				Nationality: response.Nationality,
				DateOfBirth: dateOfBirth,
				Pep:         response.Pep,
			},
		},
	}, nil
}
