package entity

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
)

type Person struct {
	DocumentNumber            string                     `json:"document_number"`
	Name                      string                     `json:"name"`
	PersonType                values.PersonType          `json:"person_type"`
	Email                     string                     `json:"email"`
	PartnerID                 string                     `json:"partner_id"`
	OfferType                 string                     `json:"offer_type"`
	ProfileID                 uuid.UUID                  `json:"profile_id"`
	EntityID                  uuid.UUID                  `json:"entity_id"`
	EntityType                values.EntityType          `json:"entity_type"`
	RoleType                  values.RoleType            `json:"role_type"`
	Individual                *Individual                `json:"individual,omitempty"`
	Company                   *Company                   `json:"company,omitempty"`
	EnrichedInformation       *EnrichedInformation       `json:"enriched_information,omitempty"`
	BlacklistStatus           *BlacklistStatus           `json:"blacklist_status,omitempty"`
	Watchlist                 *Watchlist                 `json:"watchlist,omitempty"`
	Addresses                 []Address                  `json:"addresses,omitempty"`
	Contacts                  []Contact                  `json:"contacts,omitempty"`
	NotificationRecipients    []NotificationRecipient    `json:"notification_recipients,omitempty"`
	Documents                 []Document                 `json:"documents,omitempty"`
	DocumentFiles             []DocumentFile             `json:"document_files,omitempty"`
	PEPInformation            *PepInformation            `json:"pep_information,omitempty"`
	Overrides                 []Override                 `json:"overrides,omitempty"`
	ValidationSteps           []RuleValidatorStep        `json:"validation_steps,omitempty"`
	CadastralValidationConfig *CadastralValidationConfig `json:"cadastral_validation_config,omitempty"`
}

func (person Person) HasCadastralValidationConfig() bool {
	return person.CadastralValidationConfig != nil
}

func (person Person) HasProductConfig() bool {
	return person.HasCadastralValidationConfig() && person.CadastralValidationConfig.HasProductConfig()
}

func (person Person) IsIndividual() bool {
	return person.PersonType == values.PersonTypeIndividual
}

func (person Person) IsCompany() bool {
	return person.PersonType == values.PersonTypeCompany
}

func (person Person) IsValidWatchlistIndividual() bool {
	return person.IsIndividual() &&
		person.Individual != nil &&
		person.Individual.DateOfBirth != nil
}

func (person Person) IsValidWatchlistCompany() bool {
	return person.IsCompany() &&
		person.Company != nil &&
		person.Company.LegalName != ""
}

func (person Person) ShouldCreateInternalAccount() bool {
	return person.HasCadastralValidationConfig() &&
		person.CadastralValidationConfig.HasInternalAccountCreation()
}

func (person Person) ShouldIntegrateTree() bool {
	return person.HasCadastralValidationConfig() &&
		person.CadastralValidationConfig.HasTreeIntegration()
}

func (person Person) ShouldIntegrateLimit() bool {
	return person.HasCadastralValidationConfig() &&
		person.CadastralValidationConfig.HasLimitIntegration()
}

func (person Person) ShouldEnrichProfileData() bool {
	return person.HasCadastralValidationConfig() &&
		person.CadastralValidationConfig.HasBureauEnrichment()
}

func (person Person) ShouldValidateBlacklist() bool {
	return person.HasCadastralValidationConfig() &&
		person.CadastralValidationConfig.ValidationSteps.HaveBlacklistValidation()
}

func (person Person) ShouldValidateBureau() bool {
	return person.HasCadastralValidationConfig() &&
		person.CadastralValidationConfig.ValidationSteps.HaveBureauValidation()
}

func (person Person) ShouldValidateActivityRisk() bool {
	return person.HasCadastralValidationConfig() &&
		person.CadastralValidationConfig.ValidationSteps.HaveActivityRiskValidation()
}

func (person Person) ShouldValidateOCR() bool {
	return person.HasCadastralValidationConfig() &&
		person.CadastralValidationConfig.ValidationSteps.HaveORCValidation()
}

func (person Person) ShouldValidateWatchlist() bool {
	return person.HasCadastralValidationConfig() &&
		person.CadastralValidationConfig.ValidationSteps.HaveWatchlistValidation()
}

func (person Person) ShouldValidatePEP() bool {
	return person.HasCadastralValidationConfig() &&
		person.CadastralValidationConfig.ValidationSteps.HavePEPValidation()
}

func (person Person) ShouldValidateDocuments() bool {
	return person.HasCadastralValidationConfig() &&
		person.CadastralValidationConfig.ValidationSteps.HaveDocumentsValidation()
}

func (person Person) ShouldValidateAddress() bool {
	return person.HasCadastralValidationConfig() &&
		person.CadastralValidationConfig.ValidationSteps.HaveAddressValidation()
}

func (person Person) ShouldValidateCAF() bool {
	return person.HasCadastralValidationConfig() &&
		person.CadastralValidationConfig.ValidationSteps.HaveCAFValidation()
}

func (person Person) ShouldGetDocuments() bool {
	return person.ShouldValidateDocuments() ||
		person.ShouldValidateOCR()
}

func (person Person) ShouldGetAddresses() bool {
	return person.ShouldValidateAddress() ||
		person.ShouldIntegrateTree()
}

func (person Person) ShouldGetContacts() bool {
	return person.ShouldIntegrateTree()
}

func (person Person) ShouldGetNotificationRecipients() bool {
	return person.ShouldIntegrateTree()
}

func (person Person) ShouldGetWatchlist() bool {
	return person.ShouldValidateWatchlist() ||
		person.ShouldValidatePEP()
}

func (person Person) ShouldGetBureauInformation() bool {
	return person.ShouldValidateBureau() ||
		person.ShouldValidateActivityRisk() ||
		person.ShouldEnrichProfileData()
}

func (person Person) ShouldBeEnriched() bool {
	return person.ShouldValidateCAF()
}

func (person Person) SatisfiesLegalNatureCondition(conditionValues []string) bool {
	if person.IsCompany() && person.Company != nil {
		for _, conditionValue := range conditionValues {
			if person.Company.LegalNature == conditionValue {
				return true
			}
		}
	}
	return false
}

func (person Person) SatisfiesCondition(condition Condition) bool {
	switch condition.FieldName {
	case values.LegalNatureFieldName:
		return person.SatisfiesLegalNatureCondition(condition.Values)
	default:
		return false
	}
}

func (person Person) SatisfiesAllConditions(conditions []Condition) bool {
	for _, condition := range conditions {
		if !person.SatisfiesCondition(condition) {
			return false
		}
	}
	return true
}

func (person Person) FindDocumentsEquals(documentRequired DocumentRequired) Documents {
	var documents Documents

	for _, document := range person.Documents {
		if document.IsSameTypeOf(documentRequired) && document.IsSameSubtypeOf(documentRequired) ||
			document.IsSameTypeOf(documentRequired) && !documentRequired.HasSubtype() {
			documents = append(documents, document)
		}
	}
	return documents
}

func (person Person) ShouldGetPepInformation() bool {
	return person.HasCadastralValidationConfig() &&
		person.CadastralValidationConfig.ValidationSteps.HavePEPValidation()
}
