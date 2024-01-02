package complianceProfile

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type complianceProfileService struct {
	complianceProfileRepository interfaces.ComplianceProfileRepository
	personRepository            interfaces.PersonRepository
}

func NewComplianceProfileService(profileRepository interfaces.ComplianceProfileRepository, personRepository interfaces.PersonRepository) interfaces.ComplianceProfileService {
	return &complianceProfileService{
		complianceProfileRepository: profileRepository,
		personRepository:            personRepository,
	}
}

func (ref *complianceProfileService) Get(profileID uuid.UUID) (*entity.Profile, error) {
	complianceProfile, err := ref.complianceProfileRepository.Get(profileID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if complianceProfile == nil {
		return nil, nil
	}

	if complianceProfile.LegalRepresentatives != nil {
		legalRepresentatives := make([]entity.LegalRepresentative, 0)
		for _, legalRepresentative := range complianceProfile.LegalRepresentatives {
			person, _ := ref.personRepository.Get(legalRepresentative.EntityID)
			if person != nil {
				legalRepresentative.Person = *person
			}
			legalRepresentatives = append(legalRepresentatives, legalRepresentative)
		}
		complianceProfile.LegalRepresentatives = legalRepresentatives
	}

	if complianceProfile.OwnershipStructure != nil && complianceProfile.OwnershipStructure.Shareholders != nil {
		shareholders := make([]entity.Shareholder, 0)
		for _, shareholder := range complianceProfile.OwnershipStructure.Shareholders {
			person, _ := ref.personRepository.Get(shareholder.Person.EntityID)
			if person != nil {
				shareholder.Person = *person
			}
			shareholders = append(shareholders, shareholder)
		}
		complianceProfile.OwnershipStructure.Shareholders = shareholders
	}

	if len(complianceProfile.BoardOfDirectors) > 0 {
		boardOfDirectors := make([]entity.Director, 0)
		for _, director := range complianceProfile.BoardOfDirectors {
			person, _ := ref.personRepository.Get(director.DirectorID)
			if person != nil {
				director.Person = *person
			}
			boardOfDirectors = append(boardOfDirectors, director)
		}
		complianceProfile.BoardOfDirectors = boardOfDirectors
	}

	if complianceProfile.EnrichedInformation != nil && len(complianceProfile.EnrichedInformation.EnrichedCompany.BoardOfDirectors) > 0 {
		boardOfDirectors := make([]entity.Director, 0)
		for _, director := range complianceProfile.EnrichedInformation.EnrichedCompany.BoardOfDirectors {
			person, _ := ref.personRepository.Get(director.DirectorID)
			if person != nil {
				director.Person = *person
			}
			boardOfDirectors = append(boardOfDirectors, director)
		}
		complianceProfile.EnrichedInformation.EnrichedCompany.BoardOfDirectors = boardOfDirectors
	}

	if complianceProfile.EnrichedInformation != nil && complianceProfile.EnrichedInformation.EnrichedCompany.OwnershipStructure != nil &&
		complianceProfile.EnrichedInformation.EnrichedCompany.OwnershipStructure.Shareholders != nil {
		shareholders := make([]entity.Shareholder, 0)
		for _, shareholder := range complianceProfile.EnrichedInformation.EnrichedCompany.OwnershipStructure.Shareholders {
			person, _ := ref.personRepository.Get(shareholder.Person.EntityID)
			if person != nil {
				shareholder.Person = *person
			}
			shareholders = append(shareholders, shareholder)
		}
		complianceProfile.EnrichedInformation.EnrichedCompany.OwnershipStructure.Shareholders = shareholders
	}

	return complianceProfile, err
}
