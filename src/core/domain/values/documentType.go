package values

import (
	"fmt"
)

type DocumentType = string

const (
	DocumentTypeIdentification                       DocumentType = "IDENTIFICATION"
	DocumentTypePassport                             DocumentType = "PASSPORT"
	DocumentTypeAccountOpeningContract               DocumentType = "ACCOUNT_OPENING_CONTRACT"
	DocumentTypeMandatoryStatementsAgreementEvidence DocumentType = "MANDATORY_STATEMENTS_AGREEMENT_EVIDENCE"
	DocumentTypeRegistrationForm                     DocumentType = "REGISTRATION_FORM"
	DocumentTypeCorporateDocument                    DocumentType = "CORPORATE_DOCUMENT"
	DocumentTypeProofOfAddress                       DocumentType = "PROOF_OF_ADDRESS"
	DocumentTypeSupplierAgreement                    DocumentType = "SUPPLIER_AGREEMENT"
	DocumentTypeBusinessLicense                      DocumentType = "BUSINESS_LICENSE"
	DocumentTypeFinancialStatement                   DocumentType = "FINANCIAL_STATEMENT"
	DocumentTypeInvoice                              DocumentType = "INVOICE"
	DocumentTypeAppointmentDocument                  DocumentType = "APPOINTMENT_DOCUMENT"
	DocumentTypeConstitutionDocument                 DocumentType = "CONSTITUTION_DOCUMENT"
)

var ValidDocumentTypes = map[string]DocumentType{
	string(DocumentTypeIdentification):                       DocumentTypeIdentification,
	string(DocumentTypePassport):                             DocumentTypePassport,
	string(DocumentTypeAccountOpeningContract):               DocumentTypeAccountOpeningContract,
	string(DocumentTypeMandatoryStatementsAgreementEvidence): DocumentTypeMandatoryStatementsAgreementEvidence,
	string(DocumentTypeRegistrationForm):                     DocumentTypeRegistrationForm,
	string(DocumentTypeCorporateDocument):                    DocumentTypeCorporateDocument,
	string(DocumentTypeProofOfAddress):                       DocumentTypeProofOfAddress,
	string(DocumentTypeSupplierAgreement):                    DocumentTypeSupplierAgreement,
	string(DocumentTypeBusinessLicense):                      DocumentTypeBusinessLicense,
	string(DocumentTypeFinancialStatement):                   DocumentTypeFinancialStatement,
	string(DocumentTypeInvoice):                              DocumentTypeInvoice,
	string(DocumentTypeAppointmentDocument):                  DocumentTypeAppointmentDocument,
	string(DocumentTypeConstitutionDocument):                 DocumentTypeConstitutionDocument,
}

func ParseToDocumentType(value string) (DocumentType, error) {
	if _, exists := ValidDocumentTypes[value]; !exists {
		return "", NewErrorValidation(fmt.Sprintf("%s is an invalid document type", value))
	}

	return DocumentType(value), nil
}

type DocumentSubType string

const (
	DocumentSubTypeRg                DocumentSubType = "RG"
	DocumentSubTypeRne               DocumentSubType = "RNE"
	DocumentSubTypeRnm               DocumentSubType = "RNM"
	DocumentSubTypeCnh               DocumentSubType = "CNH"
	DocumentSubTypePassport          DocumentSubType = "PASSPORT"
	DocumentSubTypeBalanceSheet      DocumentSubType = "BALANCE_SHEET"
	DocumentSubTypeBillingReport     DocumentSubType = "BILLING_REPORT"
	DocumentSubTypeMinutesOfElection DocumentSubType = "MINUTES_OF_ELECTION"
	DocumentSubTypeLetterOfAttorney  DocumentSubType = "LETTER_OF_ATTORNEY"
	DocumentSubTypeSocialContract    DocumentSubType = "SOCIAL_CONTRACT"
	DocumentSubTypeStatuteSocial     DocumentSubType = "STATUTE_SOCIAL"
	DocumentSubTypeOrganogram        DocumentSubType = "ORGANOGRAM"
)

var validDocumentSubTypes = map[string]DocumentSubType{
	string(DocumentSubTypeRg):                DocumentSubTypeRg,
	string(DocumentSubTypeRne):               DocumentSubTypeRne,
	string(DocumentSubTypeRnm):               DocumentSubTypeRnm,
	string(DocumentSubTypeCnh):               DocumentSubTypeCnh,
	string(DocumentSubTypePassport):          DocumentSubTypePassport,
	string(DocumentSubTypeBalanceSheet):      DocumentSubTypeBalanceSheet,
	string(DocumentSubTypeBillingReport):     DocumentSubTypeBillingReport,
	string(DocumentSubTypeMinutesOfElection): DocumentSubTypeMinutesOfElection,
	string(DocumentSubTypeLetterOfAttorney):  DocumentSubTypeLetterOfAttorney,
	string(DocumentSubTypeSocialContract):    DocumentSubTypeSocialContract,
	string(DocumentSubTypeStatuteSocial):     DocumentSubTypeStatuteSocial,
	string(DocumentSubTypeOrganogram):        DocumentSubTypeOrganogram,
}

func (ref *DocumentSubType) Validate() error {
	value := string(*ref)
	if _, exists := validDocumentSubTypes[value]; !exists {
		return NewErrorValidation(fmt.Sprintf("%s is an invalid document sub type", value))
	}

	return nil
}
