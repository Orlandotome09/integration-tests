package entity

import (
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DocumentsRequiredHaveAnyPendingOnApproval(t *testing.T) {
	documentsRequired := DocumentsRequired{
		{PendingOnApproval: false},
		{PendingOnApproval: false},
		{PendingOnApproval: false},
		{PendingOnApproval: true},
	}

	result := documentsRequired.HaveAnyPendingOnApproval()

	assert.True(t, result)
}

func Test_DocumentsRequired_DoesNotHaveAnyPendingOnApproval(t *testing.T) {
	documentsRequired := DocumentsRequired{
		{PendingOnApproval: false},
		{PendingOnApproval: false},
		{PendingOnApproval: false},
		{PendingOnApproval: false},
	}

	result := documentsRequired.HaveAnyPendingOnApproval()

	assert.False(t, result)
}

func Test_FilterDocumentsRequiredByConditionsSatisfied(t *testing.T) {
	documentRequiredIdentification := DocumentRequired{
		DocumentType: values2.DocumentTypeIdentification,
		Conditions: []Condition{
			{
				FieldName: values2.LegalNatureFieldName,
				Values:    []string{"1111"},
			},
		},
	}
	documentRequiredCNH := DocumentRequired{
		DocumentType: values2.DocumentTypePassport,
		Conditions: []Condition{
			{
				FieldName: values2.LegalNatureFieldName,
				Values:    []string{"2222"},
			},
		},
	}
	documentRequiredAppointment := DocumentRequired{
		DocumentType: values2.DocumentTypeAppointmentDocument,
		Conditions: []Condition{
			{
				FieldName: values2.LegalNatureFieldName,
				Values:    []string{"1111"},
			},
		},
	}
	documentsRequired := DocumentsRequired{documentRequiredIdentification, documentRequiredCNH, documentRequiredAppointment}

	person := Person{
		PersonType: values2.PersonTypeCompany,
		Company:    &Company{LegalNature: "1111"},
	}

	expected := DocumentsRequired{documentRequiredIdentification, documentRequiredAppointment}

	received := documentsRequired.FilterByConditionsSatisfied(person)

	assert.Equal(t, expected, received)
}
