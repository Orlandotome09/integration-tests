package incompleteContractRule

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	documentAdapter     *mocks.DocumentAdapter
	documentFileAdapter *mocks.DocumentFileAdapter
	incompleteRule      entity.Rule
)

func TestMain(m *testing.M) {
	documentAdapter = &mocks.DocumentAdapter{}
	documentFileAdapter = &mocks.DocumentFileAdapter{}
	os.Exit(m.Run())
}

func TestAnalyze_ContractApproved(t *testing.T) {
	contractID := uuid.New()
	DocumentID := uuid.New()
	profileID := uuid.New()
	contract := entity.Contract{ProfileID: &profileID, ContractID: &contractID, DocumentID: &DocumentID}
	incompleteRule = New(contract, documentAdapter, documentFileAdapter)

	document := &entity.Document{DocumentID: contractID, EntityID: profileID}
	documentFiles := []entity.DocumentFile{{DocumentID: contractID}}

	documentAdapter.On("GetByID", DocumentID.String()).Return(document, nil)
	documentFileAdapter.On("FindByDocumentID", DocumentID).Return(documentFiles, nil)

	expectedResult1 := entity.RuleResultV2{
		RuleSet:  values.RuleSetIncompleteContract,
		RuleName: values.RuleNameDocumentNotFound,
		Result:   values.ResultStatusApproved,
		Pending:  false,
		Metadata: nil,
	}

	expectedResult2 := entity.RuleResultV2{
		RuleSet:  values.RuleSetIncompleteContract,
		RuleName: values.RuleNameFileNotFound,
		Result:   values.ResultStatusApproved,
		Pending:  true,
		Metadata: nil,
	}

	results, err := incompleteRule.Analyze()

	require.Contains(t, results, expectedResult1)
	require.Contains(t, results, expectedResult2)
	assert.Nil(t, err)
}

func TestAnalyze_ContractRejected(t *testing.T) {
	contractID := uuid.New()
	documentID := uuid.New()
	profileID := uuid.New()
	contract := entity.Contract{ProfileID: &profileID, ContractID: &contractID, DocumentID: &documentID}
	incompleteRule = New(contract, documentAdapter, documentFileAdapter)

	document := &entity.Document{DocumentID: contractID, EntityID: uuid.New()}
	documentFiles := []entity.DocumentFile{{DocumentID: contractID}}

	documentAdapter.On("GetByID", documentID.String()).Return(document, nil)
	documentFileAdapter.On("FindByDocumentID", documentID).Return(documentFiles, nil)

	metadata, _ := json.Marshal(fmt.Sprintf("Invoice Document is associated to another profile"))
	result := entity.RuleResultV2{
		RuleSet:  values.RuleSetIncompleteContract,
		RuleName: values.RuleNameDocumentNotFound,
		Result:   values.ResultStatusRejected,
		Pending:  false,
		Metadata: metadata,
		Problems: []entity.Problem{{Code: values.ProblemCodeInvoiceAssociatedToAnotherProfile, Detail: documentID.String()}},
	}

	results, err := incompleteRule.Analyze()

	require.Contains(t, results, result)
	assert.Nil(t, err)
}

func TestAnalyze_InvoiceDocumentNotFound(t *testing.T) {
	contractID := uuid.New()
	documentID := uuid.New()
	profileID := uuid.New()
	contract := entity.Contract{ProfileID: &profileID, ContractID: &contractID, DocumentID: &documentID}
	incompleteRule = New(contract, documentAdapter, documentFileAdapter)

	var document *entity.Document = nil

	documentAdapter.On("GetByID", documentID.String()).Return(document, nil)

	metadata, _ := json.Marshal(fmt.Sprintf("Invoice Document: %v not found", documentID))
	result := entity.RuleResultV2{
		RuleSet:  values.RuleSetIncompleteContract,
		RuleName: values.RuleNameDocumentNotFound,
		Result:   values.ResultStatusIncomplete,
		Pending:  false,
		Metadata: metadata,
		Problems: []entity.Problem{{Code: values.ProblemCodeInvoiceDocumentNotFound, Detail: documentID.String()}},
	}

	results, err := incompleteRule.Analyze()

	require.Contains(t, results, result)
	assert.Nil(t, err)
}

func TestAnalyze_InvoiceFileNotFound(t *testing.T) {
	contractID := uuid.New()
	documentID := uuid.New()
	profileID := uuid.New()
	contract := entity.Contract{ProfileID: &profileID, ContractID: &contractID, DocumentID: &documentID}
	incompleteRule = New(contract, documentAdapter, documentFileAdapter)

	document := &entity.Document{DocumentID: contractID, EntityID: profileID}
	documentFiles := []entity.DocumentFile{}

	documentAdapter.On("GetByID", documentID.String()).Return(document, nil)
	documentFileAdapter.On("FindByDocumentID", documentID).Return(documentFiles, nil)

	metadata, _ := json.Marshal(fmt.Sprintf("Invoice File not found for Document: %v",
		documentID))
	result := entity.RuleResultV2{
		RuleSet:  values.RuleSetIncompleteContract,
		RuleName: values.RuleNameFileNotFound,
		Result:   values.ResultStatusIncomplete,
		Pending:  false,
		Metadata: metadata,
		Problems: []entity.Problem{{Code: values.ProblemCodeInvoiceFileNotFound, Detail: documentID.String()}},
	}

	results, err := incompleteRule.Analyze()

	require.Contains(t, results, result)
	assert.Nil(t, err)
}

func TestGetName(t *testing.T) {
	id := uuid.MustParse("4d2d0cc3-f173-4a78-86c8-5b9de2fbf8d2")
	contract := entity.Contract{ContractID: &id}
	incompleteRule = New(contract, documentAdapter, documentFileAdapter)

	receivedName := incompleteRule.Name()

	expected := values.RuleSetIncompleteContract

	assert.Equal(t, expected, receivedName)
}

func TestAnalyze_ContractInvoiceRequired(t *testing.T) {
	contractID := uuid.New()
	profileID := uuid.New()
	contract := entity.Contract{ProfileID: &profileID, ContractID: &contractID, DocumentID: nil}
	incompleteRule = New(contract, documentAdapter, documentFileAdapter)

	document := &entity.Document{DocumentID: contractID, EntityID: uuid.New()}
	documentFiles := []entity.DocumentFile{{DocumentID: contractID}}

	documentAdapter.On("GetByID", mock.Anything).Return(document, nil)
	documentFileAdapter.On("FindByDocumentID", mock.Anything).Return(documentFiles, nil)

	metadata, _ := json.Marshal(fmt.Sprintf("Invoice is required"))
	result := entity.RuleResultV2{
		RuleSet:  values.RuleSetIncompleteContract,
		RuleName: values.RuleNameDocumentNotFound,
		Result:   values.ResultStatusIncomplete,
		Pending:  false,
		Metadata: metadata,
		Problems: []entity.Problem{{Code: values.ProblemCodeInvoiceIsRequired, Detail: ""}},
	}

	results, err := incompleteRule.Analyze()

	require.Contains(t, results, result)
	assert.Nil(t, err)
}
