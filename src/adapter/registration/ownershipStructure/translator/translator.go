package translator

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/ownershipStructure/http/contracts"
	"gopkg.in/errgo.v2/fmt/errors"
)

type OwnershipStructureTranslator interface {
	Translate(response contracts.OwnershipStructureResponse) (*entity.OwnershipStructure, error)
}

type ownershipStructureTranslator struct{}

func New() OwnershipStructureTranslator {
	return &ownershipStructureTranslator{}
}

func (ref *ownershipStructureTranslator) Translate(response contracts.OwnershipStructureResponse) (*entity.OwnershipStructure, error) {

	shareholders, err := translateShareholders(response.Shareholders)
	if err != nil {
		return nil, err
	}

	return &entity.OwnershipStructure{
		FinalBeneficiariesCount: response.FinalBeneficiariesCounted,
		ShareholdingSum:         response.ShareholdingSum,
		Shareholders:            shareholders,
	}, nil
}

func translateShareholders(response contracts.Shareholders) ([]entity.Shareholder, error) {

	shareholders := make([]entity.Shareholder, 0)

	for _, shareholder := range response {

		domainShareholder, err := translateShareholder(shareholder)
		if err != nil {
			return nil, err
		}

		shareholders = append(shareholders, *domainShareholder)

		domainShareholders, err := translateShareholders(shareholder.Shareholders)
		if err != nil {
			return nil, err
		}

		shareholders = append(shareholders, domainShareholders...)

	}
	return removeDuplicated(shareholders), nil
}

func translateShareholder(response contracts.Shareholder) (*entity.Shareholder, error) {

	var birthDate *time.Time = nil

	if response.BirthDate != "" {
		result, err := time.Parse("02/01/2006", response.BirthDate)
		if err != nil {
			return nil, errors.Newf("[OwnershipStructureTranslator] Error translating birthDate %v", birthDate)
		}
		birthDate = &result
	}

	shareholder := &entity.Shareholder{
		ShareholderID:    &response.ShareholderID,
		OwnershipPercent: response.Shareholding,
		Person: entity.Person{
			DocumentNumber: response.DocumentNumber,
			Name:           response.Name,
			PersonType:     response.Type,
			EntityID:       response.ShareholderID,
			EntityType:     values.EntityTypeShareholder,
			RoleType:       values.RoleTypeShareholder,
		},
	}

	if shareholder.Person.PersonType == values.PersonTypeIndividual {
		shareholder.Person.Individual = &entity.Individual{
			DateOfBirth:         birthDate,
			DateOfBirthInputted: birthDate,
			Pep:                 &response.Pep,
		}
	}

	if shareholder.Person.PersonType == values.PersonTypeCompany {
		shareholder.Person.Company = &entity.Company{
			LegalName: response.Name,
		}
	}

	return shareholder, nil
}

func removeDuplicated(shareholders []entity.Shareholder) []entity.Shareholder {
	keys := make(map[string]bool)
	var filtered []entity.Shareholder

	for _, shareholder := range shareholders {
		_, exists := keys[shareholder.Person.DocumentNumber]
		if !exists {
			keys[shareholder.Person.DocumentNumber] = true
			filtered = append(filtered, shareholder)
		}
	}
	return filtered
}
