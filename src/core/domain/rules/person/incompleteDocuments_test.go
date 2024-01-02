package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewIncompleteDocumentsAnalyzer(t *testing.T) {
	analyzerDocuments := NewIncompleteDocumentsAnalyzer([]entity.DocumentRequired{})

	assert.NotNil(t, analyzerDocuments)
}

func TestIncompleteDocumentsAnalyzer_Analyze_Success(t *testing.T) {
	documentsRequired := []entity.DocumentRequired{
		{
			DocumentType:      values.DocumentTypeInvoice,
			FileRequired:      true,
			PendingOnApproval: true,
			Conditions:        []entity.Condition{},
		},
	}
	analyzerDocuments := NewIncompleteDocumentsAnalyzer(documentsRequired)

	profileID := uuid.New()
	documentID := uuid.New()
	documentFileID := uuid.New()
	person := entity.Person{
		ProfileID: profileID,
		EntityID:  profileID,
		Documents: []entity.Document{{
			DocumentType: values.DocumentTypeInvoice,
			DocumentID:   documentID,
			EntityID:     profileID,
		}},
		DocumentFiles: []entity.DocumentFile{{
			DocumentFileID: &documentFileID,
			DocumentID:     documentID,
		}},
	}

	result, err := analyzerDocuments.Analyze(person)

	assert.Nil(t, err)
	assert.Equal(t, values.ResultStatusApproved, result.Result)
	assert.True(t, result.Pending)
}

func TestIncompleteDocumentsAnalyzer_Analyze_Subtype_Success(t *testing.T) {
	documentsRequired := []entity.DocumentRequired{
		{
			DocumentType:      values.DocumentTypeIdentification,
			FileRequired:      true,
			PendingOnApproval: true,
			Conditions:        []entity.Condition{},
		},
	}
	analyzerDocuments := NewIncompleteDocumentsAnalyzer(documentsRequired)

	profileID := uuid.New()
	documentID := uuid.New()
	documentFileID := uuid.New()
	person := entity.Person{
		ProfileID: profileID,
		EntityID:  profileID,
		Documents: []entity.Document{{
			DocumentType:    values.DocumentTypeIdentification,
			DocumentSubType: string(values.DocumentSubTypeRg),
			DocumentID:      documentID,
			EntityID:        profileID,
		}},
		DocumentFiles: []entity.DocumentFile{{
			DocumentFileID: &documentFileID,
			DocumentID:     documentID,
		}},
	}

	result, err := analyzerDocuments.Analyze(person)

	assert.Nil(t, err)
	assert.Equal(t, values.ResultStatusApproved, result.Result)
	assert.True(t, result.Pending)
}

func TestIncompleteDocumentsAnalyzer_Analyze_Subtype_Pending(t *testing.T) {
	documentsRequired := []entity.DocumentRequired{
		{
			DocumentType:      values.DocumentTypeIdentification,
			DocumentSubType:   string(values.DocumentSubTypeRg),
			FileRequired:      true,
			PendingOnApproval: true,
			Conditions:        []entity.Condition{},
		},
	}
	analyzerDocuments := NewIncompleteDocumentsAnalyzer(documentsRequired)

	profileID := uuid.New()
	documentID := uuid.New()
	documentFileID := uuid.New()
	person := entity.Person{
		ProfileID: profileID,
		EntityID:  profileID,
		Documents: []entity.Document{{
			DocumentType: values.DocumentTypeIdentification,
			DocumentID:   documentID,
			EntityID:     profileID,
		}},
		DocumentFiles: []entity.DocumentFile{{
			DocumentFileID: &documentFileID,
			DocumentID:     documentID,
		}},
	}

	result, err := analyzerDocuments.Analyze(person)

	assert.Nil(t, err)
	assert.Equal(t, values.ResultStatusIncomplete, result.Result)
	assert.Equal(t, values.ProblemCodeDocumentNotFoundRG, result.Problems[0].Code)
}

func TestIncompleteDocumentsAnalyzer_Analyze_Pending(t *testing.T) {
	documentsRequired := []entity.DocumentRequired{
		{
			DocumentType:      values.DocumentTypeInvoice,
			FileRequired:      true,
			PendingOnApproval: true,
			Conditions:        []entity.Condition{},
		}, {
			DocumentType:      values.DocumentTypeBusinessLicense,
			FileRequired:      false,
			PendingOnApproval: false,
			Conditions:        []entity.Condition{},
		},
	}
	analyzerDocuments := NewIncompleteDocumentsAnalyzer(documentsRequired)

	profileID := uuid.New()
	documentID := uuid.New()
	person := entity.Person{
		ProfileID: profileID,
		EntityID:  profileID,
		Documents: []entity.Document{{
			DocumentType: values.DocumentTypeInvoice,
			DocumentID:   documentID,
			EntityID:     profileID,
		}},
	}

	result, err := analyzerDocuments.Analyze(person)

	assert.Nil(t, err)
	assert.Equal(t, values.ResultStatusIncomplete, result.Result)
	assert.Equal(t, values.ProblemCodeFileNotFoundInvoice, result.Problems[0].Code)
	assert.Equal(t, values.ProblemCodeDocumentNotFoundBusinessLicense, result.Problems[1].Code)
}
