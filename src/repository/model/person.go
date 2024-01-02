package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"

	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type Person struct {
	EntityID                  uuid.UUID `gorm:"primaryKey;type:uuid"`
	DocumentNumber            string
	Name                      string
	PersonType                string
	Email                     string
	PartnerID                 string
	OfferType                 string
	ProfileID                 uuid.UUID `gorm:"type:uuid"`
	EntityType                string
	RoleType                  string
	Individual                postgres.Jsonb `gorm:"type:jsonb"`
	Company                   postgres.Jsonb `gorm:"type:jsonb"`
	EnrichmentInformation     postgres.Jsonb `gorm:"type:jsonb"`
	BlacklistStatus           postgres.Jsonb `gorm:"type:jsonb"`
	Watchlist                 postgres.Jsonb `gorm:"type:jsonb"`
	Addresses                 postgres.Jsonb `gorm:"type:jsonb"`
	Documents                 postgres.Jsonb `gorm:"type:jsonb"`
	DocumentFiles             postgres.Jsonb `gorm:"type:jsonb"`
	Overrides                 postgres.Jsonb `gorm:"type:jsonb"`
	CadastralValidationConfig postgres.Jsonb `gorm:"type:jsonb"`
	NotificationRecipients    postgres.Jsonb `gorm:"type:jsonb"`
	PEPInformation            postgres.Jsonb `gorm:"type:jsonb"`
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

func NewPersonFromDomain(domainPerson entity2.Person) (*Person, error) {

	individual, err := marshal(domainPerson.Individual)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	company, err := marshal(domainPerson.Company)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	enrichmentInformation, err := marshal(domainPerson.EnrichedInformation)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	blacklist, err := marshal(domainPerson.BlacklistStatus)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	watchlist, err := marshal(domainPerson.Watchlist)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	addresses, err := marshal(domainPerson.Addresses)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	documents, err := marshal(domainPerson.Documents)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	documentFiles, err := marshal(domainPerson.DocumentFiles)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	overrides, err := marshal(domainPerson.Overrides)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cadastralValidationConfig, err := marshal(domainPerson.CadastralValidationConfig)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	notificationRecipients, err := marshal(domainPerson.NotificationRecipients)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pepInformation, err := marshal(domainPerson.PEPInformation)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Person{
		EntityID:                  domainPerson.EntityID,
		DocumentNumber:            domainPerson.DocumentNumber,
		Name:                      domainPerson.Name,
		PersonType:                domainPerson.PersonType,
		Email:                     domainPerson.Email,
		PartnerID:                 domainPerson.PartnerID,
		OfferType:                 domainPerson.OfferType,
		ProfileID:                 domainPerson.ProfileID,
		EntityType:                domainPerson.EntityType.ToString(),
		RoleType:                  domainPerson.RoleType,
		Individual:                postgres.Jsonb{RawMessage: individual},
		Company:                   postgres.Jsonb{RawMessage: company},
		EnrichmentInformation:     postgres.Jsonb{RawMessage: enrichmentInformation},
		BlacklistStatus:           postgres.Jsonb{RawMessage: blacklist},
		Watchlist:                 postgres.Jsonb{RawMessage: watchlist},
		Addresses:                 postgres.Jsonb{RawMessage: addresses},
		Documents:                 postgres.Jsonb{RawMessage: documents},
		DocumentFiles:             postgres.Jsonb{RawMessage: documentFiles},
		Overrides:                 postgres.Jsonb{RawMessage: overrides},
		CadastralValidationConfig: postgres.Jsonb{RawMessage: cadastralValidationConfig},
		NotificationRecipients:    postgres.Jsonb{RawMessage: notificationRecipients},
		PEPInformation:            postgres.Jsonb{RawMessage: pepInformation},
	}, nil
}

func (ref *Person) ToDomain() (*entity2.Person, error) {

	var individual *entity2.Individual
	err := unmarshal(ref.Individual.RawMessage, &individual)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var company *entity2.Company
	err = unmarshal(ref.Company.RawMessage, &company)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var enrichedInformation *entity2.EnrichedInformation
	err = unmarshal(ref.EnrichmentInformation.RawMessage, &enrichedInformation)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var blacklist *entity2.BlacklistStatus
	err = unmarshal(ref.BlacklistStatus.RawMessage, &blacklist)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var watchlist *entity2.Watchlist
	err = unmarshal(ref.Watchlist.RawMessage, &watchlist)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var addresses []entity2.Address
	err = unmarshal(ref.Addresses.RawMessage, &addresses)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var documents []entity2.Document
	err = unmarshal(ref.Documents.RawMessage, &documents)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var documentFiles []entity2.DocumentFile
	err = unmarshal(ref.DocumentFiles.RawMessage, &documentFiles)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var overrides []entity2.Override
	err = unmarshal(ref.Overrides.RawMessage, &overrides)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var cadastralValidationConfig *entity2.CadastralValidationConfig
	err = unmarshal(ref.CadastralValidationConfig.RawMessage, &cadastralValidationConfig)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var notificationRecipients []entity2.NotificationRecipient
	err = unmarshal(ref.NotificationRecipients.RawMessage, &notificationRecipients)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var pepInformation *entity2.PepInformation
	err = unmarshal(ref.PEPInformation.RawMessage, &pepInformation)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &entity2.Person{
		DocumentNumber:            ref.DocumentNumber,
		Name:                      ref.Name,
		PersonType:                ref.PersonType,
		Email:                     ref.Email,
		PartnerID:                 ref.PartnerID,
		OfferType:                 ref.OfferType,
		ProfileID:                 ref.ProfileID,
		EntityID:                  ref.EntityID,
		EntityType:                values.EntityType(ref.EntityType),
		RoleType:                  ref.RoleType,
		Individual:                individual,
		Company:                   company,
		EnrichedInformation:       enrichedInformation,
		BlacklistStatus:           blacklist,
		Watchlist:                 watchlist,
		Addresses:                 addresses,
		Documents:                 documents,
		DocumentFiles:             documentFiles,
		Overrides:                 overrides,
		CadastralValidationConfig: cadastralValidationConfig,
		NotificationRecipients:    notificationRecipients,
		PEPInformation:            pepInformation,
	}, nil

}
