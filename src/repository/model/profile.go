package model

import (
	"encoding/json"
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
)

type ComplianceProfile struct {
	ProfileID            uuid.UUID  `gorm:"primaryKey;type:uuid"`
	ParentID             *uuid.UUID `gorm:"type:uuid"`
	LegacyID             string
	CallbackUrl          string
	LegalRepresentatives postgres.Jsonb `gorm:"type:jsonb"`
	OwnershipStructure   postgres.Jsonb `gorm:"type:jsonb"`
	BoardOfDirectors     postgres.Jsonb `gorm:"type:jsonb"`
	ExpirationDate       *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Person               postgres.Jsonb `gorm:"type:jsonb"`
}

func NewProfileFromDomain(profile entity.Profile) (*ComplianceProfile, error) {

	legalRepresentatives, err := marshal(profile.LegalRepresentatives)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ownershipStructure, err := marshal(profile.OwnershipStructure)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	boardOfDirectors, err := marshal(profile.BoardOfDirectors)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	person, err := marshal(profile.Person)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &ComplianceProfile{
		ProfileID:            *profile.ProfileID,
		ParentID:             profile.ParentID,
		LegacyID:             profile.LegacyID,
		CallbackUrl:          profile.CallbackUrl,
		LegalRepresentatives: postgres.Jsonb{RawMessage: legalRepresentatives},
		OwnershipStructure:   postgres.Jsonb{RawMessage: ownershipStructure},
		BoardOfDirectors:     postgres.Jsonb{RawMessage: boardOfDirectors},
		ExpirationDate:       profile.ExpirationDate,
		CreatedAt:            profile.CreatedAt,
		UpdatedAt:            profile.UpdatedAt,
		Person:               postgres.Jsonb{RawMessage: person},
	}, nil
}

func (ref *ComplianceProfile) ToDomain() (*entity.Profile, error) {

	var legalRepresentatives []entity.LegalRepresentative
	if err := unmarshal(ref.LegalRepresentatives.RawMessage, &legalRepresentatives); err != nil {
		return nil, errors.WithStack(err)
	}

	var ownershipStructure *entity.OwnershipStructure
	if err := unmarshal(ref.OwnershipStructure.RawMessage, &ownershipStructure); err != nil {
		return nil, errors.WithStack(err)
	}

	var boardOfDirectors []entity.Director
	if err := unmarshal(ref.BoardOfDirectors.RawMessage, &boardOfDirectors); err != nil {
		return nil, errors.WithStack(err)
	}

	var person entity.Person
	if err := unmarshal(ref.Person.RawMessage, &person); err != nil {
		return nil, errors.WithStack(err)
	}

	return &entity.Profile{
		Person:               person,
		ProfileID:            &ref.ProfileID,
		ParentID:             ref.ParentID,
		LegacyID:             ref.LegacyID,
		CallbackUrl:          ref.CallbackUrl,
		LegalRepresentatives: legalRepresentatives,
		OwnershipStructure:   ownershipStructure,
		BoardOfDirectors:     boardOfDirectors,
		ExpirationDate:       ref.ExpirationDate,
		CreatedAt:            ref.CreatedAt,
		UpdatedAt:            ref.UpdatedAt,
	}, nil
}

func marshal(entity interface{}) ([]byte, error) {

	if entity == nil {
		return nil, nil
	}

	return json.Marshal(entity)
}

func unmarshal(jsonMessage []byte, entity interface{}) error {

	if jsonMessage == nil {
		entity = nil
		return nil
	}

	return json.Unmarshal(jsonMessage, entity)
}
