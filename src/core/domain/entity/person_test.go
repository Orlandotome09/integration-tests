package entity

import (
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_PersonHasCatalog(t *testing.T) {
	person := Person{CadastralValidationConfig: &CadastralValidationConfig{}}

	hasCatalog := person.HasCadastralValidationConfig()

	assert.True(t, hasCatalog)
}

func Test_PersonHasNotCatalog(t *testing.T) {
	person := Person{CadastralValidationConfig: nil}

	hasCatalog := person.HasCadastralValidationConfig()

	assert.False(t, hasCatalog)
}

func Test_PersonHasProductConfig(t *testing.T) {
	person := Person{CadastralValidationConfig: &CadastralValidationConfig{ProductConfig: &ProductConfig{}}}

	hasProductConfig := person.HasProductConfig()

	assert.True(t, hasProductConfig)
}

func Test_PersonDoesNotHaveProductConfig(t *testing.T) {
	person := Person{CadastralValidationConfig: &CadastralValidationConfig{ProductConfig: nil}}

	hasProductConfig := person.HasProductConfig()

	assert.False(t, hasProductConfig)
}

func Test_PersonIsValidWatchlistIndividual(t *testing.T) {
	dateOfBirth := time.Now()
	person := Person{
		PersonType: values2.PersonTypeIndividual,
		Individual: &Individual{DateOfBirth: &dateOfBirth},
	}

	result := person.IsValidWatchlistIndividual()

	assert.True(t, result)
}

func Test_PersonIsValidWatchlistCompany(t *testing.T) {
	person := Person{
		PersonType: values2.PersonTypeCompany,
		Company:    &Company{LegalName: "legal name"},
	}

	result := person.IsValidWatchlistCompany()

	assert.True(t, result)
}

func Test_PersonShouldCreateInternalAccount(t *testing.T) {
	person := Person{CadastralValidationConfig: &CadastralValidationConfig{ProductConfig: &ProductConfig{CreateBexsAccount: true}}}

	result := person.ShouldCreateInternalAccount()

	assert.True(t, result)
}

func Test_PersonShouldNotCreateInternalAccount(t *testing.T) {
	person := Person{CadastralValidationConfig: &CadastralValidationConfig{ProductConfig: &ProductConfig{CreateBexsAccount: false}}}

	result := person.ShouldCreateInternalAccount()

	assert.False(t, result)
}

func Test_PersonShouldIntegrateTree(t *testing.T) {
	person := Person{CadastralValidationConfig: &CadastralValidationConfig{ProductConfig: &ProductConfig{TreeIntegration: true}}}

	result := person.ShouldIntegrateTree()

	assert.True(t, result)
}

func Test_PersonShouldNotIntegrateTree(t *testing.T) {
	person := Person{CadastralValidationConfig: &CadastralValidationConfig{ProductConfig: &ProductConfig{TreeIntegration: false}}}

	result := person.ShouldIntegrateTree()

	assert.False(t, result)
}

func Test_PersonShouldIntegrateLimit(t *testing.T) {
	person := Person{CadastralValidationConfig: &CadastralValidationConfig{ProductConfig: &ProductConfig{LimitIntegration: true}}}

	result := person.ShouldIntegrateLimit()

	assert.True(t, result)
}

func Test_PersonShouldNotIntegrateLimit(t *testing.T) {
	person := Person{CadastralValidationConfig: &CadastralValidationConfig{ProductConfig: &ProductConfig{LimitIntegration: false}}}

	result := person.ShouldIntegrateLimit()

	assert.False(t, result)
}

func Test_PersonShouldEnrichProfileData(t *testing.T) {
	person := Person{CadastralValidationConfig: &CadastralValidationConfig{ProductConfig: &ProductConfig{EnrichProfileWithBureauData: true}}}

	result := person.ShouldEnrichProfileData()

	assert.True(t, result)
}

func Test_PersonShouldNotEnrichProfileData(t *testing.T) {
	person := Person{CadastralValidationConfig: &CadastralValidationConfig{ProductConfig: &ProductConfig{EnrichProfileWithBureauData: false}}}

	result := person.ShouldEnrichProfileData()

	assert.False(t, result)
}

func Test_PersonShouldValidateBlackList(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{BlackListParams: &BlackListParams{}},
				},
			},
		},
	}

	result := person.ShouldValidateBlacklist()

	assert.True(t, result)
}

func Test_PersonShouldNotValidateBlackList(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{BlackListParams: nil},
				},
			},
		},
	}

	result := person.ShouldValidateBlacklist()

	assert.False(t, result)
}

func Test_PersonShouldValidateBureau(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{BureauParams: &BureauParams{}},
				},
			},
		},
	}

	result := person.ShouldValidateBureau()

	assert.True(t, result)
}

func Test_PersonShouldNotValidateBureau(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{BureauParams: nil},
				},
			},
		},
	}

	result := person.ShouldValidateBureau()

	assert.False(t, result)
}

func Test_PersonShouldValidateActivityRisk(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{ActivityRiskParams: &ActivityRiskParams{}},
				},
			},
		},
	}

	result := person.ShouldValidateActivityRisk()

	assert.True(t, result)
}

func Test_PersonShouldNotValidateActivityRisk(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{ActivityRiskParams: nil},
				},
			},
		},
	}

	result := person.ShouldValidateActivityRisk()

	assert.False(t, result)
}

func Test_PersonShouldValidateOCR(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{DOAParams: &DOAParams{}},
				},
			},
		},
	}

	result := person.ShouldValidateOCR()

	assert.True(t, result)
}

func Test_PersonShouldNotValidateOCR(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{DOAParams: nil},
				},
			},
		},
	}

	result := person.ShouldValidateOCR()

	assert.False(t, result)
}

func Test_PersonShouldValidateWatchlist(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{WatchListParams: &WatchListParams{}},
				},
			},
		},
	}

	result := person.ShouldValidateWatchlist()

	assert.True(t, result)
}

func Test_PersonShouldNotValidateWatchlist(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{WatchListParams: nil},
				},
			},
		},
	}

	result := person.ShouldValidateWatchlist()

	assert.False(t, result)
}

func Test_PersonShouldValidatePEP(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{PepParams: &PepParams{}},
				},
			},
		},
	}

	result := person.ShouldValidatePEP()

	assert.True(t, result)
}

func Test_PersonShouldNotValidatePEP(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{PepParams: nil},
				},
			},
		},
	}

	result := person.ShouldValidatePEP()

	assert.False(t, result)
}

func Test_PersonShouldValidateDocuments(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{
						IncompleteParams: &IncompleteParams{
							DocumentsRequired: []DocumentRequired{{DocumentType: values2.DocumentTypeIdentification}}}},
				},
			},
		},
	}

	result := person.ShouldValidateDocuments()

	assert.True(t, result)
}

func Test_PersonShouldNotValidateDocuments(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{IncompleteParams: &IncompleteParams{DocumentsRequired: nil}},
				},
			},
		},
	}

	result := person.ShouldValidateDocuments()

	assert.False(t, result)
}

func Test_PersonShouldValidateAddress(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{
						IncompleteParams: &IncompleteParams{
							AddressRequired: true}},
				},
			},
		},
	}

	result := person.ShouldValidateAddress()

	assert.True(t, result)
}

func Test_PersonShouldNotValidateAddress(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{IncompleteParams: &IncompleteParams{AddressRequired: false}},
				},
			},
		},
	}

	result := person.ShouldValidateAddress()

	assert.False(t, result)
}

func Test_ShouldGetDocumentsWhenShouldValidateDocuments(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{
						IncompleteParams: &IncompleteParams{
							DocumentsRequired: []DocumentRequired{{DocumentType: values2.DocumentTypeIdentification}}}},
				},
			},
		},
	}

	result := person.ShouldGetDocuments()

	assert.True(t, result)
}
func Test_ShouldGetDocumentsWhenShouldValidateOCR(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{DOAParams: &DOAParams{}},
				},
			},
		},
	}

	result := person.ShouldGetDocuments()

	assert.True(t, result)
}

func Test_ShouldGetAddressesWhenShouldValidateAddress(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{
						IncompleteParams: &IncompleteParams{
							AddressRequired: true}},
				},
			},
		},
	}

	result := person.ShouldGetAddresses()

	assert.True(t, result)
}

func Test_ShouldGetAddressesWhenShouldIntegrateTree(t *testing.T) {
	person := Person{CadastralValidationConfig: &CadastralValidationConfig{ProductConfig: &ProductConfig{TreeIntegration: true}}}

	result := person.ShouldGetAddresses()

	assert.True(t, result)
}

func Test_ShouldGetWatchlistWhenShouldValidateWatchlist(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{WatchListParams: &WatchListParams{}},
				},
			},
		},
	}

	result := person.ShouldGetWatchlist()

	assert.True(t, result)
}

func Test_ShouldGetWatchlistWhenShouldValidatePEP(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{PepParams: &PepParams{}},
				},
			},
		},
	}

	result := person.ShouldGetWatchlist()

	assert.True(t, result)
}

func Test_ShouldGetBureauInformationWhenShouldValidateBureau(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{BureauParams: &BureauParams{}},
				},
			},
		},
	}

	result := person.ShouldGetBureauInformation()

	assert.True(t, result)
}

func Test_ShouldGetBureauInformationWhenShouldValidateActivityRisk(t *testing.T) {
	person := Person{
		CadastralValidationConfig: &CadastralValidationConfig{
			ValidationSteps: ValidationSteps{
				{
					RulesConfig: &RuleSetConfig{ActivityRiskParams: &ActivityRiskParams{}},
				},
			},
		},
	}

	result := person.ShouldGetBureauInformation()

	assert.True(t, result)
}

func Test_ShouldGetBureauInformationWhenShouldEnrichProfileData(t *testing.T) {
	person := Person{CadastralValidationConfig: &CadastralValidationConfig{ProductConfig: &ProductConfig{EnrichProfileWithBureauData: true}}}

	result := person.ShouldGetBureauInformation()

	assert.True(t, result)
}

func Test_PersonSatisfiesNatureCondition(t *testing.T) {
	conditionValues := []string{"1111", "2222", "3333"}
	person := Person{
		PersonType: values2.PersonTypeCompany,
		Company:    &Company{LegalNature: "1111"},
	}

	result := person.SatisfiesLegalNatureCondition(conditionValues)

	assert.True(t, result)
}

func Test_PersonDoesNotSatisfyLegalNatureCondition(t *testing.T) {
	conditionValues := []string{"1111", "2222", "3333"}
	person := Person{
		PersonType: values2.PersonTypeCompany,
		Company:    &Company{LegalNature: "4444"},
	}

	result := person.SatisfiesLegalNatureCondition(conditionValues)

	assert.False(t, result)
}

func Test_PersonSatisfiesCondition(t *testing.T) {
	condition := Condition{
		FieldName: values2.LegalNatureFieldName,
		Values:    []string{"1111", "2222", "3333"},
	}
	person := Person{
		PersonType: values2.PersonTypeCompany,
		Company:    &Company{LegalNature: "3333"},
	}

	result := person.SatisfiesCondition(condition)

	assert.True(t, result)
}

func Test_PersonDoesNotSatisfyCondition(t *testing.T) {
	condition := Condition{
		FieldName: values2.LegalNatureFieldName,
		Values:    []string{"1111", "2222", "3333"},
	}
	person := Person{
		PersonType: values2.PersonTypeCompany,
		Company:    &Company{LegalNature: "0000"},
	}

	result := person.SatisfiesCondition(condition)

	assert.False(t, result)
}

func Test_PersonSatisfiesAllConditions(t *testing.T) {
	conditions := []Condition{
		{
			FieldName: values2.LegalNatureFieldName,
			Values:    []string{"1111", "2222", "3333"},
		},
	}
	person := Person{
		PersonType: values2.PersonTypeCompany,
		Company:    &Company{LegalNature: "3333"},
	}

	result := person.SatisfiesAllConditions(conditions)

	assert.True(t, result)
}

func Test_PersonSatisfiesAllConditions_whenNoConditionsToBeValidated(t *testing.T) {
	var conditions []Condition
	person := Person{
		EntityID: uuid.New(),
	}

	result := person.SatisfiesAllConditions(conditions)

	assert.True(t, result)
}

func Test_PersonDoesNotSatisfyAllConditions(t *testing.T) {
	conditions := []Condition{
		{
			FieldName: values2.LegalNatureFieldName,
			Values:    []string{"1111", "2222", "3333"},
		},
		{
			FieldName: values2.FieldName("XXX"),
			Values:    []string{"1111", "2222", "3333"},
		},
	}
	person := Person{
		PersonType: values2.PersonTypeCompany,
		Company:    &Company{LegalNature: "3333"},
	}

	result := person.SatisfiesAllConditions(conditions)

	assert.False(t, result)
}
