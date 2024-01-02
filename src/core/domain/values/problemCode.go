package values

import (
	"fmt"
)

type ProblemCode = string

const (
	ProblemCodeInvalidProblemCode                                   ProblemCode = "INVALID_PROBLEM_CODE"
	ProblemCodeInvoiceIsRequired                                    ProblemCode = "INVOICE_IS_REQUIRED"
	ProblemCodeInvoiceDocumentNotFound                              ProblemCode = "INVOICE_DOCUMENT_NOT_FOUND"
	ProblemCodeInvoiceAssociatedToAnotherProfile                    ProblemCode = "INVOICE_ASSOCIATED_TO_ANOTHER_PROFILE"
	ProblemCodeInvoiceFileNotFound                                  ProblemCode = "INVOICE_FILE_NOT_FOUND"
	ProblemCodeStateNotFound                                        ProblemCode = "STATE_NOT_FOUND"
	ProblemCodeProfileNotApproved                                   ProblemCode = "PROFILE_NOT_APPROVED"
	ProblemCodeLegalRepresentativeNotApproved                       ProblemCode = "LEGAL_REPRESENTATIVE_NOT_APPROVED"
	ProblemCodeShareholdingNotAchieveMinimumRequired                ProblemCode = "SHAREHOLDING_NOT_ACHIEVE_MINIMUM_REQUIRED"
	ProblemCodeShareholderNotApproved                               ProblemCode = "SHAREHOLDER_NOT_APPROVED"
	ProblemCodeNotFoundAtBureau                                     ProblemCode = "BUREAU_NOT_FOUND"
	ProblemCodeBureauStatusNotRegular                               ProblemCode = "BUREAU_STATUS_NOT_REGULAR"
	ProblemCodeDateOfBirthRequired                                  ProblemCode = "DATE_OF_BIRTH_REQUIRED"
	ProblemCodeInputtedOrEnrichedDateOfBirthRequired                ProblemCode = "INPUTTED_OR_ENRICHED_DATE_OF_BIRTH_REQUIRED"
	ProblemCodePhoneRequired                                        ProblemCode = "PHONE_REQUIRED"
	ProblemCodeEmailRequired                                        ProblemCode = "EMAIL_REQUIRED"
	ProblemCodePepRequired                                          ProblemCode = "PEP_REQUIRED"
	ProblemCodeLastNameRequired                                     ProblemCode = "LAST_NAME_REQUIRED"
	ProblemCodeDocumentNotFoundInvoice                              ProblemCode = "DOCUMENT_NOT_FOUND_INVOICE"
	ProblemCodeDocumentNotFoundIdentification                       ProblemCode = "DOCUMENT_NOT_FOUND_IDENTIFICATION"
	ProblemCodeDocumentNotFoundRG                                   ProblemCode = "DOCUMENT_NOT_FOUND_RG"
	ProblemCodeDocumentNotFoundRNE                                  ProblemCode = "DOCUMENT_NOT_FOUND_RNE"
	ProblemCodeDocumentNotFoundRNM                                  ProblemCode = "DOCUMENT_NOT_FOUND_RNM"
	ProblemCodeDocumentNotFoundCNH                                  ProblemCode = "DOCUMENT_NOT_FOUND_CNH"
	ProblemCodeDocumentNotFoundPassport                             ProblemCode = "DOCUMENT_NOT_FOUND_PASSPORT"
	ProblemCodeDocumentNotFoundAccountOpeningContract               ProblemCode = "DOCUMENT_NOT_FOUND_ACCOUNT_OPENING_CONTRACT"
	ProblemCodeDocumentNotFoundMandatoryStatementsAgreementEvidence ProblemCode = "DOCUMENT_NOT_FOUND_MANDATORY_STATEMENTS_AGREEMENT_EVIDENCE"
	ProblemCodeDocumentNotFoundRegistrationForm                     ProblemCode = "DOCUMENT_NOT_FOUND_REGISTRATION_FORM"
	ProblemCodeDocumentNotFoundCorporateDocument                    ProblemCode = "DOCUMENT_NOT_FOUND_CORPORATE_DOCUMENT"
	ProblemCodeDocumentNotFoundOrganogram                           ProblemCode = "DOCUMENT_NOT_FOUND_ORGANOGRAM"
	ProblemCodeDocumentNotFoundProofOfAddress                       ProblemCode = "DOCUMENT_NOT_FOUND_PROOF_OF_ADDRESS"
	ProblemCodeDocumentNotFoundSupplierAgreement                    ProblemCode = "DOCUMENT_NOT_FOUND_SUPPLIER_AGREEMENT"
	ProblemCodeDocumentNotFoundBusinessLicense                      ProblemCode = "DOCUMENT_NOT_FOUND_BUSINESS_LICENSE"
	ProblemCodeDocumentNotFoundFinancialStatement                   ProblemCode = "DOCUMENT_NOT_FOUND_FINANCIAL_STATEMENT"
	ProblemCodeDocumentNotFoundBalanceSheet                         ProblemCode = "DOCUMENT_NOT_FOUND_BALANCE_SHEET"
	ProblemCodeDocumentNotFoundBillingReport                        ProblemCode = "DOCUMENT_NOT_FOUND_BILLING_REPORT"
	ProblemCodeDocumentNotFoundProofOfFinancialStanding             ProblemCode = "DOCUMENT_NOT_FOUND_PROOF_OF_FINANCIAL_STANDING"
	ProblemCodeDocumentNotFoundAppointmentDocument                  ProblemCode = "DOCUMENT_NOT_FOUND_APPOINTMENT_DOCUMENT"
	ProblemCodeDocumentNotFoundMinutesOfElection                    ProblemCode = "DOCUMENT_NOT_FOUND_MINUTES_OF_ELECTION"
	ProblemCodeDocumentNotFoundLetterOfAttorney                     ProblemCode = "DOCUMENT_NOT_FOUND_LETTER_OF_ATTORNEY"
	ProblemCodeDocumentNotFoundConstitutionDocument                 ProblemCode = "DOCUMENT_NOT_FOUND_CONSTITUTION_DOCUMENT"
	ProblemCodeDocumentNotFoundStatuteSocial                        ProblemCode = "DOCUMENT_NOT_FOUND_STATUTE_SOCIAL"
	ProblemCodeDocumentNotFoundSocialContract                       ProblemCode = "DOCUMENT_NOT_FOUND_SOCIAL_CONTRACT"
	ProblemCodeDocumentNotProofOfShareholderChain                   ProblemCode = "DOCUMENT_NOT_FOUND_PROOF_OF_SHAREHOLDER_CHAIN"
	ProblemCodeFileNotFoundInvoice                                  ProblemCode = "FILE_NOT_FOUND_INVOICE"
	ProblemCodeFileNotFoundIdentification                           ProblemCode = "FILE_NOT_FOUND_IDENTIFICATION"
	ProblemCodeFileNotFoundRG                                       ProblemCode = "FILE_NOT_FOUND_RG"
	ProblemCodeFileNotFoundRNE                                      ProblemCode = "FILE_NOT_FOUND_RNE"
	ProblemCodeFileNotFoundRNM                                      ProblemCode = "FILE_NOT_FOUND_RNM"
	ProblemCodeFileNotFoundCNH                                      ProblemCode = "FILE_NOT_FOUND_CNH"
	ProblemCodeFileNotFoundPassport                                 ProblemCode = "FILE_NOT_FOUND_PASSPORT"
	ProblemCodeFileNotFoundAccountOpeningContract                   ProblemCode = "FILE_NOT_FOUND_ACCOUNT_OPENING_CONTRACT"
	ProblemCodeFileNotFoundMandatoryStatementsAgreementEvidence     ProblemCode = "FILE_NOT_FOUND_MANDATORY_STATEMENTS_AGREEMENT_EVIDENCE"
	ProblemCodeFileNotFoundRegistrationForm                         ProblemCode = "FILE_NOT_FOUND_REGISTRATION_FORM"
	ProblemCodeFileNotFoundCorporateDocument                        ProblemCode = "FILE_NOT_FOUND_CORPORATE_DOCUMENT"
	ProblemCodeFileNotFoundOrganogram                               ProblemCode = "FILE_NOT_FOUND_ORGANOGRAM"
	ProblemCodeFileNotFoundProofOfAddress                           ProblemCode = "FILE_NOT_FOUND_PROOF_OF_ADDRESS"
	ProblemCodeFileNotFoundSupplierAgreement                        ProblemCode = "FILE_NOT_FOUND_SUPPLIER_AGREEMENT"
	ProblemCodeFileNotFoundBusinessLicense                          ProblemCode = "FILE_NOT_FOUND_BUSINESS_LICENSE"
	ProblemCodeFileNotFoundFinancialStatement                       ProblemCode = "FILE_NOT_FOUND_FINANCIAL_STATEMENT"
	ProblemCodeFileNotFoundBalanceSheet                             ProblemCode = "FILE_NOT_FOUND_BALANCE_SHEET"
	ProblemCodeFileNotFoundBillingReport                            ProblemCode = "FILE_NOT_FOUND_BILLING_REPORT"
	ProblemCodeFileNotFoundProofOfFinancialStanding                 ProblemCode = "FILE_NOT_FOUND_PROOF_OF_FINANCIAL_STANDING"
	ProblemCodeFileNotFoundAppointmentDocument                      ProblemCode = "FILE_NOT_FOUND_APPOINTMENT_DOCUMENT"
	ProblemCodeFileNotFoundMinutesOfElection                        ProblemCode = "FILE_NOT_FOUND_MINUTES_OF_ELECTION"
	ProblemCodeFileNotFoundLetterOfAttorney                         ProblemCode = "FILE_NOT_FOUND_LETTER_OF_ATTORNEY"
	ProblemCodeFileNotFoundConstitutionDocument                     ProblemCode = "FILE_NOT_FOUND_CONSTITUTION_DOCUMENT"
	ProblemCodeFileNotFoundStatuteSocial                            ProblemCode = "FILE_NOT_FOUND_STATUTE_SOCIAL"
	ProblemCodeFileNotFoundSocialContract                           ProblemCode = "FILE_NOT_FOUND_SOCIAL_CONTRACT"
	ProblemCodeFileNotFoundProofOfShareholderChain                  ProblemCode = "FILE_NOT_FOUND_PROOF_OF_SHAREHOLDER_CHAIN"
	ProblemCodeAddressNotFound                                      ProblemCode = "ADDRESS_NOT_FOUND"
	ProblemCodeIdentificationDocumentNotFound                       ProblemCode = "IDENTIFICATION_DOCUMENT_NOT_FOUND"
	ProblemCodeFrontFileNotFound                                    ProblemCode = "FRONT_FILE_NOT_FOUND"
	ProblemCodeBackFileNotFound                                     ProblemCode = "BACK_FILE_NOT_FOUND"
	ProblemCodeEconomicalActivityRiskUndefined                      ProblemCode = "ECONOMICAL_ACTIVITY_RISK_UNDEFINED"
	ProblemCodeEconomicalActivityRiskHigh                           ProblemCode = "ECONOMICAL_ACTIVITY_RISK_HIGH"
	ProblemCodeDateOfBirthNotInputtedOrEnriched                     ProblemCode = "DATE_OF_BIRTH_NOT_INPUTTED_OR_ENRICHED"
	ProblemCodePersonFoundOnWatchlist                               ProblemCode = "PERSON_FOUND_ON_WATCHLIST"
	ProblemCodeDirectorNotApproved                                  ProblemCode = "DIRECTOR_NOT_APPROVED"
	ProblemCodeBoardOfDirectorsIncomplete                           ProblemCode = "BOARD_OF_DIRECTORS_INCOMPLETE"
	ProblemCodeNotFoundCafAnalysis                                  ProblemCode = "NOT_FOUND_CAF_ANALYSIS"
	ProblemCodeNotFoundEnrichedInformation                          ProblemCode = "NOT_FOUND_ENRICHED_INFORMATION"
	ProblemCodeCafAnalysisPending                                   ProblemCode = "CAF_ANALYSIS_PENDING"
	ProblemCodePersonHasInsufficientMinimumIncome                   ProblemCode = "PERSON_HAS_INSUFFICIENT_INCOME"
	ProblemCodeCompanyHasInsufficientBilling                        ProblemCode = "COMPANY_HAS_INSUFFICIENT_BILLING"
)

var validProblemCodes = map[string]ProblemCode{
	ProblemCodeInvoiceIsRequired:                                    ProblemCodeInvoiceIsRequired,
	ProblemCodeInvoiceDocumentNotFound:                              ProblemCodeInvoiceDocumentNotFound,
	ProblemCodeInvoiceAssociatedToAnotherProfile:                    ProblemCodeInvoiceAssociatedToAnotherProfile,
	ProblemCodeInvoiceFileNotFound:                                  ProblemCodeInvoiceFileNotFound,
	ProblemCodeStateNotFound:                                        ProblemCodeStateNotFound,
	ProblemCodeProfileNotApproved:                                   ProblemCodeProfileNotApproved,
	ProblemCodeLegalRepresentativeNotApproved:                       ProblemCodeLegalRepresentativeNotApproved,
	ProblemCodeShareholdingNotAchieveMinimumRequired:                ProblemCodeShareholdingNotAchieveMinimumRequired,
	ProblemCodeShareholderNotApproved:                               ProblemCodeShareholderNotApproved,
	ProblemCodeNotFoundAtBureau:                                     ProblemCodeNotFoundAtBureau,
	ProblemCodeBureauStatusNotRegular:                               ProblemCodeBureauStatusNotRegular,
	ProblemCodeDateOfBirthRequired:                                  ProblemCodeDateOfBirthRequired,
	ProblemCodeInputtedOrEnrichedDateOfBirthRequired:                ProblemCodeInputtedOrEnrichedDateOfBirthRequired,
	ProblemCodePhoneRequired:                                        ProblemCodePhoneRequired,
	ProblemCodeEmailRequired:                                        ProblemCodeEmailRequired,
	ProblemCodePepRequired:                                          ProblemCodePepRequired,
	ProblemCodeLastNameRequired:                                     ProblemCodeLastNameRequired,
	ProblemCodeDocumentNotFoundInvoice:                              ProblemCodeDocumentNotFoundInvoice,
	ProblemCodeDocumentNotFoundIdentification:                       ProblemCodeDocumentNotFoundIdentification,
	ProblemCodeDocumentNotFoundRG:                                   ProblemCodeDocumentNotFoundRG,
	ProblemCodeDocumentNotFoundRNE:                                  ProblemCodeDocumentNotFoundRNE,
	ProblemCodeDocumentNotFoundRNM:                                  ProblemCodeDocumentNotFoundRNM,
	ProblemCodeDocumentNotFoundCNH:                                  ProblemCodeDocumentNotFoundCNH,
	ProblemCodeDocumentNotFoundPassport:                             ProblemCodeDocumentNotFoundPassport,
	ProblemCodeDocumentNotFoundAccountOpeningContract:               ProblemCodeDocumentNotFoundAccountOpeningContract,
	ProblemCodeDocumentNotFoundMandatoryStatementsAgreementEvidence: ProblemCodeDocumentNotFoundMandatoryStatementsAgreementEvidence,
	ProblemCodeDocumentNotFoundRegistrationForm:                     ProblemCodeDocumentNotFoundRegistrationForm,
	ProblemCodeDocumentNotFoundCorporateDocument:                    ProblemCodeDocumentNotFoundCorporateDocument,
	ProblemCodeDocumentNotFoundOrganogram:                           ProblemCodeDocumentNotFoundOrganogram,
	ProblemCodeDocumentNotFoundProofOfAddress:                       ProblemCodeDocumentNotFoundProofOfAddress,
	ProblemCodeDocumentNotFoundSupplierAgreement:                    ProblemCodeDocumentNotFoundSupplierAgreement,
	ProblemCodeDocumentNotFoundBusinessLicense:                      ProblemCodeDocumentNotFoundBusinessLicense,
	ProblemCodeDocumentNotFoundFinancialStatement:                   ProblemCodeDocumentNotFoundFinancialStatement,
	ProblemCodeDocumentNotFoundBalanceSheet:                         ProblemCodeDocumentNotFoundBalanceSheet,
	ProblemCodeDocumentNotFoundBillingReport:                        ProblemCodeDocumentNotFoundBillingReport,
	ProblemCodeDocumentNotFoundProofOfFinancialStanding:             ProblemCodeDocumentNotFoundProofOfFinancialStanding,
	ProblemCodeDocumentNotFoundAppointmentDocument:                  ProblemCodeDocumentNotFoundAppointmentDocument,
	ProblemCodeDocumentNotFoundMinutesOfElection:                    ProblemCodeDocumentNotFoundMinutesOfElection,
	ProblemCodeDocumentNotFoundLetterOfAttorney:                     ProblemCodeDocumentNotFoundLetterOfAttorney,
	ProblemCodeDocumentNotFoundConstitutionDocument:                 ProblemCodeDocumentNotFoundConstitutionDocument,
	ProblemCodeDocumentNotFoundStatuteSocial:                        ProblemCodeDocumentNotFoundStatuteSocial,
	ProblemCodeDocumentNotFoundSocialContract:                       ProblemCodeDocumentNotFoundSocialContract,
	ProblemCodeDocumentNotProofOfShareholderChain:                   ProblemCodeDocumentNotProofOfShareholderChain,
	ProblemCodeFileNotFoundInvoice:                                  ProblemCodeFileNotFoundInvoice,
	ProblemCodeFileNotFoundIdentification:                           ProblemCodeFileNotFoundIdentification,
	ProblemCodeFileNotFoundRG:                                       ProblemCodeFileNotFoundRG,
	ProblemCodeFileNotFoundRNE:                                      ProblemCodeFileNotFoundRNE,
	ProblemCodeFileNotFoundRNM:                                      ProblemCodeFileNotFoundRNM,
	ProblemCodeFileNotFoundCNH:                                      ProblemCodeFileNotFoundCNH,
	ProblemCodeFileNotFoundPassport:                                 ProblemCodeFileNotFoundPassport,
	ProblemCodeFileNotFoundAccountOpeningContract:                   ProblemCodeFileNotFoundAccountOpeningContract,
	ProblemCodeFileNotFoundMandatoryStatementsAgreementEvidence:     ProblemCodeFileNotFoundMandatoryStatementsAgreementEvidence,
	ProblemCodeFileNotFoundRegistrationForm:                         ProblemCodeFileNotFoundRegistrationForm,
	ProblemCodeFileNotFoundCorporateDocument:                        ProblemCodeFileNotFoundCorporateDocument,
	ProblemCodeFileNotFoundOrganogram:                               ProblemCodeFileNotFoundOrganogram,
	ProblemCodeFileNotFoundProofOfAddress:                           ProblemCodeFileNotFoundProofOfAddress,
	ProblemCodeFileNotFoundSupplierAgreement:                        ProblemCodeFileNotFoundSupplierAgreement,
	ProblemCodeFileNotFoundBusinessLicense:                          ProblemCodeFileNotFoundBusinessLicense,
	ProblemCodeFileNotFoundFinancialStatement:                       ProblemCodeFileNotFoundFinancialStatement,
	ProblemCodeFileNotFoundBalanceSheet:                             ProblemCodeFileNotFoundBalanceSheet,
	ProblemCodeFileNotFoundBillingReport:                            ProblemCodeFileNotFoundBillingReport,
	ProblemCodeFileNotFoundProofOfFinancialStanding:                 ProblemCodeFileNotFoundProofOfFinancialStanding,
	ProblemCodeFileNotFoundAppointmentDocument:                      ProblemCodeFileNotFoundAppointmentDocument,
	ProblemCodeFileNotFoundMinutesOfElection:                        ProblemCodeFileNotFoundMinutesOfElection,
	ProblemCodeFileNotFoundLetterOfAttorney:                         ProblemCodeFileNotFoundLetterOfAttorney,
	ProblemCodeFileNotFoundConstitutionDocument:                     ProblemCodeFileNotFoundConstitutionDocument,
	ProblemCodeFileNotFoundStatuteSocial:                            ProblemCodeFileNotFoundStatuteSocial,
	ProblemCodeFileNotFoundSocialContract:                           ProblemCodeFileNotFoundSocialContract,
	ProblemCodeFileNotFoundProofOfShareholderChain:                  ProblemCodeFileNotFoundProofOfShareholderChain,
	ProblemCodeAddressNotFound:                                      ProblemCodeAddressNotFound,
	ProblemCodeIdentificationDocumentNotFound:                       ProblemCodeIdentificationDocumentNotFound,
	ProblemCodeFrontFileNotFound:                                    ProblemCodeFrontFileNotFound,
	ProblemCodeBackFileNotFound:                                     ProblemCodeBackFileNotFound,
	ProblemCodeEconomicalActivityRiskUndefined:                      ProblemCodeEconomicalActivityRiskUndefined,
	ProblemCodeEconomicalActivityRiskHigh:                           ProblemCodeEconomicalActivityRiskHigh,
	ProblemCodeDateOfBirthNotInputtedOrEnriched:                     ProblemCodeDateOfBirthNotInputtedOrEnriched,
	ProblemCodePersonFoundOnWatchlist:                               ProblemCodePersonFoundOnWatchlist,
	ProblemCodeDirectorNotApproved:                                  ProblemCodeDirectorNotApproved,
	ProblemCodeBoardOfDirectorsIncomplete:                           ProblemCodeBoardOfDirectorsIncomplete,
	ProblemCodeNotFoundCafAnalysis:                                  ProblemCodeNotFoundCafAnalysis,
	ProblemCodeNotFoundEnrichedInformation:                          ProblemCodeNotFoundEnrichedInformation,
	ProblemCodeCafAnalysisPending:                                   ProblemCodeCafAnalysisPending,
	ProblemCodePersonHasInsufficientMinimumIncome:                   ProblemCodePersonHasInsufficientMinimumIncome,
	ProblemCodeCompanyHasInsufficientBilling:                        ProblemCodeCompanyHasInsufficientBilling,
}

func ProblemCodeParser(problemCode string) error {
	_, in := validProblemCodes[problemCode]
	if !in {
		return NewErrorValidation(fmt.Sprintf("%s is an invalid problem code", problemCode))
	}
	return nil
}
