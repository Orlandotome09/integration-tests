package legalRepresentativeTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/legalRepresentative/http/contracts"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type LegalRepresentativeTranslator interface {
	Translate(response contracts.LegalRepresentativeResponse) (*entity.LegalRepresentative, error)
}

type legalRepresentativeTranslator struct{}

func New() LegalRepresentativeTranslator {
	return &legalRepresentativeTranslator{}
}

func (ref *legalRepresentativeTranslator) Translate(response contracts.LegalRepresentativeResponse) (*entity.LegalRepresentative, error) {
	legalRepresentativeID, err := uuid.Parse(response.LegalRepresentativeID)
	if err != nil {
		return nil, errors.New("[LegalRepresentativeTranslator] LegalRepresentativeID is not a valid uuid " + legalRepresentativeID.String())
	}

	profileID, err := uuid.Parse(response.ProfileID)
	if err != nil {
		return nil, errors.New("[LegalRepresentativeTranslator] ProfileID is not a valid uuid " + response.ProfileID)
	}

	var birthDate *time.Time = nil

	if response.BirthDate != "" {
		result, err := time.Parse("02/01/2006", response.BirthDate)
		if err != nil {
			return nil, errors.New("[LegalRepresentativeTranslator] BirthDate is not a valid date " + response.BirthDate)
		}
		birthDate = &result
	}

	return &entity.LegalRepresentative{
		LegalRepresentativeID: legalRepresentativeID,
		ExpirationDate:        response.ExpirationDate,
		Person: entity.Person{
			DocumentNumber: response.DocumentNumber,
			Name:           response.FullName,
			PersonType:     values.PersonTypeIndividual,
			Email:          response.Email,
			PartnerID:      response.PartnerID,
			OfferType:      response.OfferType,
			ProfileID:      profileID,
			EntityID:       legalRepresentativeID,
			EntityType:     values.EntityTypeLegalRepresentative,
			RoleType:       values.RoleTypeLegalRepresentative,
			Individual: &entity.Individual{
				Nationality: response.Nationality,
				DateOfBirth: birthDate,
				Pep:         response.Pep,
				Phones: []entity.Phone{
					{
						Number: response.Phone,
					},
				},
				ForeignTaxResidency:   response.ForeignTaxResidency,
				CountryOfTaxResidency: response.CountryOfTaxResidency,
			},
		},
	}, nil
}
