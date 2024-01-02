package values

import (
	"fmt"
)

type RuleName string

const (
	// IsUnderAge Set
	RuleNameCustomerIsUnderAge RuleName = "CUSTOMER_IS_UNDER_AGE"

	// Serasa Set
	RuleNameNotFoundInSerasa    RuleName = "CUSTOMER_NOT_FOUND_IN_SERASA"
	RuleNameHasProblemsInSerasa RuleName = "CUSTOMER_HAS_PROBLEMS_IN_SERASA"

	// Watchlist Set
	RuleNameOccurrenceInWatchlist RuleName = "WATCHLIST"

	// Blacklist Set
	RuleNameOccurrenceInBlacklist RuleName = "BLACKLIST"

	// Incomplete Set
	RuleNameRequiredFieldsNotFound RuleName = "REQUIRED_FIELDS_NOT_FOUND"

	// Precondition Contract Set
	RuleNameProfileApproved RuleName = "PROFILE_APPROVED"

	// Test Set
	RuleNameTest RuleName = "TEST"

	// DOA Set
	RuleNameDOAValidation   RuleName = "DOA_VALIDATION"
	RuleNameDOAFileNotfound RuleName = "DOA_FILE_NOT_FOUND"

	// Manual Block Set
	RuleNameManualBlock RuleName = "MANUAL_BLOCK"

	// Shareholders Set
	RuleNameShareholding RuleName = "SHAREHOLDING"
	RuleNameShareholders RuleName = "SHAREHOLDERS"

	// PEP Set
	RuleNamePep RuleName = "PEP"

	// LegalRepresentative Set
	RuleNameLegalRepresentativesResult RuleName = "LEGAL_REPRESENTATIVES_RESULT"

	// Risky activity
	RuleNameHighRiskActivity RuleName = "HIGH_RISK_ACTIVITY"

	// Generics: used in any set
	RuleNameDocumentNotFound RuleName = "DOCUMENT_NOT_FOUND"
	RuleNameFileNotFound     RuleName = "FILE_NOT_FOUND"
	RuleNameAddressNotFound  RuleName = "ADDRESS_NOT_FOUND"
	RuleNameBlocked          RuleName = "BLOCKED"
	RuleNameInactive         RuleName = "INACTIVE"

	// RuleSetBoardOfDirectors
	RuleNameBoardOfDirectorsComplete RuleName = "BOARD_OF_DIRECTORS_COMPLETE"
	RuleNameBoardOfDirectorsResult   RuleName = "BOARD_OF_DIRECTORS_RESULT"

	// Manual Validation Set
	RuleNameManualValidation RuleName = "PROFILE_DATA"

	RuleNameCafAnalysis RuleName = "CAF_ANALYSIS_RESULT"

	RuleNameInsufficientBilling RuleName = "INSUFFICIENT_BILLING"
	RuleNameInsufficientIncome  RuleName = "INSUFFICIENT_INCOME"
)

var validRuleNames = map[string]RuleName{
	RuleNameCustomerIsUnderAge.ToString():         RuleNameCustomerIsUnderAge,
	RuleNameNotFoundInSerasa.ToString():           RuleNameNotFoundInSerasa,
	RuleNameHasProblemsInSerasa.ToString():        RuleNameHasProblemsInSerasa,
	RuleNameOccurrenceInWatchlist.ToString():      RuleNameOccurrenceInWatchlist,
	RuleNameOccurrenceInBlacklist.ToString():      RuleNameOccurrenceInBlacklist,
	RuleNameRequiredFieldsNotFound.ToString():     RuleNameRequiredFieldsNotFound,
	RuleNameProfileApproved.ToString():            RuleNameProfileApproved,
	RuleNameTest.ToString():                       RuleNameTest,
	RuleNameDOAValidation.ToString():              RuleNameDOAValidation,
	RuleNameDOAFileNotfound.ToString():            RuleNameDOAFileNotfound,
	RuleNameManualBlock.ToString():                RuleNameManualBlock,
	RuleNamePep.ToString():                        RuleNamePep,
	RuleNameLegalRepresentativesResult.ToString(): RuleNameLegalRepresentativesResult,
	RuleNameHighRiskActivity.ToString():           RuleNameHighRiskActivity,
	RuleNameDocumentNotFound.ToString():           RuleNameDocumentNotFound,
	RuleNameFileNotFound.ToString():               RuleNameFileNotFound,
	RuleNameAddressNotFound.ToString():            RuleNameAddressNotFound,
	RuleNameBlocked.ToString():                    RuleNameBlocked,
	RuleNameInactive.ToString():                   RuleNameInactive,
	RuleNameShareholding.ToString():               RuleNameShareholding,
	RuleNameShareholders.ToString():               RuleNameShareholders,
	RuleNameBoardOfDirectorsComplete.ToString():   RuleNameBoardOfDirectorsComplete,
	RuleNameBoardOfDirectorsResult.ToString():     RuleNameBoardOfDirectorsResult,
	RuleNameManualValidation.ToString():           RuleNameManualValidation,
	RuleNameCafAnalysis.ToString():                RuleNameCafAnalysis,
	RuleNameInsufficientBilling.ToString():        RuleNameInsufficientBilling,
	RuleNameInsufficientIncome.ToString():         RuleNameInsufficientIncome,
}

func (ruleName RuleName) Validate() error {
	_, in := validRuleNames[ruleName.ToString()]
	if !in {
		return NewErrorValidation(fmt.Sprintf("%s is an invalid rule name", ruleName))
	}
	return nil
}

func (ruleName RuleName) ToString() string {
	return string(ruleName)
}
