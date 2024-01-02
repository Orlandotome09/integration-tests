package person

import (
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/doa"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/doa/contracts"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var fileService *mocks.FileAdapter
var doaService *doa.MockDOAService
var analyzer entity2.Rule

func TestMain(m *testing.M) {
	fileService = &mocks.FileAdapter{}
	doaService = &doa.MockDOAService{}
	os.Exit(m.Run())
}

func TestAnalyze_should_start_extraction_successfully(t *testing.T) {
	profileID := uuid.New()
	documentID := uuid.New()
	frontFileID := uuid.New()
	backFileID := uuid.New()
	person := entity2.Person{
		ProfileID: profileID,
		EntityID:  profileID,
		Documents: []entity2.Document{{
			DocumentID:      documentID,
			DocumentType:    values2.DocumentTypeIdentification,
			DocumentSubType: string(values2.DocumentSubTypeRg),
			DocumentFields:  entity2.DocumentFields{Number: "888", Name: "Lola"},
		}},
		DocumentFiles: []entity2.DocumentFile{
			{DocumentFileID: &frontFileID, FileID: uuid.New(), FileSide: values2.FileSideFront, DocumentID: documentID},
			{DocumentFileID: &backFileID, FileID: uuid.New(), FileSide: values2.FileSideBack, DocumentID: documentID},
		},
	}
	analyzer = NewDoaAnalyzer(fileService, doaService, person, nil, nil)

	var lastResult *entity2.DOAResult = nil

	frontFileURI := "/url1"
	backFileURI := "/url2"

	response := &entity2.DOAExtraction{RequestID: uuid.New()}

	doaService.On("FindLastResult", &person.ProfileID, &documentID).Return(lastResult, nil)

	doaService.On("RequestExtraction", &person.DocumentFiles[0], frontFileURI, &person.DocumentFiles[1], backFileURI,
		&person.Documents[0], profileID).Return(response, nil)

	fileService.On("GetUrl", person.DocumentFiles[0].FileID).Return(frontFileURI, nil)

	fileService.On("GetUrl", person.DocumentFiles[1].FileID).Return(backFileURI, nil)

	doaResult := &entity2.DOAResult{ID: response.RequestID,
		EntityID:   profileID,
		DocumentID: person.Documents[0].DocumentID,
		FileIDs:    []uuid.UUID{person.DocumentFiles[0].FileID, person.DocumentFiles[1].FileID},
		Status:     values2.DOAStatusValidating,
	}

	doaService.On("Save", doaResult).Return(doaResult, nil)

	metadata, _ := json.Marshal("Awaiting for DOA calculate score")
	expected := []entity2.RuleResultV2{
		{
			RuleSet:  values2.RuleSetDOA,
			RuleName: values2.RuleNameDocumentNotFound,
			Result:   values2.ResultStatusApproved,
		},
		{
			RuleSet:  values2.RuleSetDOA,
			RuleName: values2.RuleNameDOAFileNotfound,
			Result:   values2.ResultStatusApproved,
		},
		{
			RuleSet:  values2.RuleSetDOA,
			RuleName: values2.RuleNameDOAValidation,
			Result:   values2.ResultStatusAnalysing,
			Metadata: metadata,
		},
	}

	received, err := analyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)

}

func TestAnalyze_documentNotFound(t *testing.T) {
	profileID := uuid.New()
	person := entity2.Person{ProfileID: profileID, EntityID: profileID}
	analyzer = NewDoaAnalyzer(fileService, doaService, person, nil, nil)

	_, err := analyzer.Analyze()

	assert.Nil(t, err)
}

func TestAnalyze_should_not_find_front_file(t *testing.T) {
	profileID := uuid.New()
	documentID := uuid.New()
	person := entity2.Person{ProfileID: profileID, EntityID: profileID,
		Documents: []entity2.Document{
			{
				DocumentID:   documentID,
				DocumentType: values2.DocumentTypeIdentification,
			},
		},
		DocumentFiles: []entity2.DocumentFile{
			{FileID: uuid.New(), FileSide: values2.FileSideBack, DocumentID: documentID},
		},
	}
	analyzer = NewDoaAnalyzer(fileService, doaService, person, nil, nil)

	metadata, _ := json.Marshal("FRONT FILE Not Found")
	expected := []entity2.RuleResultV2{
		{
			RuleSet:  values2.RuleSetDOA,
			RuleName: values2.RuleNameDocumentNotFound,
			Result:   values2.ResultStatusApproved,
		},
		{
			RuleSet:  values2.RuleSetDOA,
			RuleName: values2.RuleNameDOAFileNotfound,
			Result:   values2.ResultStatusIncomplete,
			Metadata: metadata,
			Problems: []entity2.Problem{{Code: values2.ProblemCodeFrontFileNotFound, Detail: ""}},
		},
	}

	received, err := analyzer.Analyze()

	assert.Nil(t, err)
	assert.Equal(t, expected, received)

}

func TestAnalyze_should_not_find_back_file(t *testing.T) {
	profileID := uuid.New()
	documentID := uuid.New()
	frontFileID := uuid.New()
	person := entity2.Person{ProfileID: profileID, EntityID: profileID,
		Documents: []entity2.Document{
			{
				DocumentID:   documentID,
				DocumentType: values2.DocumentTypeIdentification,
			},
		},
		DocumentFiles: []entity2.DocumentFile{
			{DocumentFileID: &frontFileID, FileID: uuid.New(), FileSide: values2.FileSideFront, DocumentID: documentID},
		},
	}
	analyzer = NewDoaAnalyzer(fileService, doaService, person, nil, nil)

	metadata, _ := json.Marshal("BACK FILE Not Found")
	expected := []entity2.RuleResultV2{
		{
			RuleSet:  values2.RuleSetDOA,
			RuleName: values2.RuleNameDocumentNotFound,
			Result:   values2.ResultStatusApproved,
		},
		{
			RuleSet:  values2.RuleSetDOA,
			RuleName: values2.RuleNameDOAFileNotfound,
			Result:   values2.ResultStatusIncomplete,
			Metadata: metadata,
			Problems: []entity2.Problem{{Code: values2.ProblemCodeBackFileNotFound, Detail: ""}},
		},
	}

	received, err := analyzer.Analyze()

	assert.Nil(t, err)

	assert.Equal(t, expected, received)
}

func TestAnalyze_errFindingLastDoaResult(t *testing.T) {
	profileID := uuid.New()
	documentID := uuid.New()
	frontFileID := uuid.New()
	backFileID := uuid.New()
	person := entity2.Person{ProfileID: profileID, EntityID: profileID,
		Documents: []entity2.Document{{
			DocumentID:      documentID,
			DocumentType:    values2.DocumentTypeIdentification,
			DocumentSubType: string(values2.DocumentSubTypeRg),
			DocumentFields:  entity2.DocumentFields{Number: "888", Name: "Lola"},
		}},
		DocumentFiles: []entity2.DocumentFile{
			{DocumentFileID: &frontFileID, FileID: uuid.New(), FileSide: values2.FileSideFront, DocumentID: documentID},
			{DocumentFileID: &backFileID, FileID: uuid.New(), FileSide: values2.FileSideBack, DocumentID: documentID},
		},
	}
	analyzer = NewDoaAnalyzer(fileService, doaService, person, nil, nil)
	var lastResult *entity2.DOAResult = nil
	errFindingLastDoaResult := errors.New("error")

	doaService.On("FindLastResult", &person.ProfileID, &documentID).Return(lastResult, errFindingLastDoaResult)

	_, err := analyzer.Analyze()

	assert.Equal(t, errFindingLastDoaResult, err)
}

func TestAnalyze_errGettingUrlForFrontFile(t *testing.T) {
	profileID := uuid.New()
	documentID := uuid.New()
	frontFileID := uuid.New()
	backFileID := uuid.New()
	person := entity2.Person{ProfileID: profileID, EntityID: profileID,
		Documents: []entity2.Document{{
			DocumentID:      documentID,
			DocumentType:    values2.DocumentTypeIdentification,
			DocumentSubType: string(values2.DocumentSubTypeRg),
			DocumentFields:  entity2.DocumentFields{Number: "888", Name: "Lola"},
		}},
		DocumentFiles: []entity2.DocumentFile{
			{DocumentFileID: &frontFileID, FileID: uuid.New(), FileSide: values2.FileSideFront, DocumentID: documentID},
			{DocumentFileID: &backFileID, FileID: uuid.New(), FileSide: values2.FileSideBack, DocumentID: documentID},
		},
	}
	analyzer = NewDoaAnalyzer(fileService, doaService, person, nil, nil)

	var lastResult *entity2.DOAResult = nil
	errGettingUrlForFrontFile := errors.New("err for frontfile")

	doaService.On("FindLastResult", &person.ProfileID, &documentID).Return(lastResult, nil)

	fileService.On("GetUrl", person.DocumentFiles[0].FileID).Return("", errGettingUrlForFrontFile)

	_, err := analyzer.Analyze()

	assert.Equal(t, errGettingUrlForFrontFile, err)
}

func TestAnalyze_errGettingUrlForBackFile(t *testing.T) {
	profileID := uuid.New()
	documentID := uuid.New()
	frontFileID := uuid.New()
	backFileID := uuid.New()
	person := entity2.Person{ProfileID: profileID, EntityID: profileID,
		Documents: []entity2.Document{{
			DocumentID:      documentID,
			DocumentType:    values2.DocumentTypeIdentification,
			DocumentSubType: string(values2.DocumentSubTypeRg),
			DocumentFields:  entity2.DocumentFields{Number: "888", Name: "Lola"},
		}},
		DocumentFiles: []entity2.DocumentFile{
			{DocumentFileID: &frontFileID, FileID: uuid.New(), FileSide: values2.FileSideFront, DocumentID: documentID},
			{DocumentFileID: &backFileID, FileID: uuid.New(), FileSide: values2.FileSideBack, DocumentID: documentID},
		},
	}
	analyzer = NewDoaAnalyzer(fileService, doaService, person, nil, nil)

	var lastResult *entity2.DOAResult = nil
	errGettingUrlForBackFile := errors.New("err for backfile")

	doaService.On("FindLastResult", &person.ProfileID, &documentID).Return(lastResult, nil)

	fileService.On("GetUrl", person.DocumentFiles[0].FileID).Return("", nil)

	fileService.On("GetUrl", person.DocumentFiles[1].FileID).Return("", errGettingUrlForBackFile)

	_, err := analyzer.Analyze()

	assert.Equal(t, errGettingUrlForBackFile, err)
}

func TestAnalyze_notChanged_success(t *testing.T) {
	profileID := uuid.New()
	documentID := uuid.New()
	frontFileID := uuid.New()
	backFileID := uuid.New()
	person := entity2.Person{ProfileID: profileID, EntityID: profileID,
		Documents: []entity2.Document{{
			DocumentID:      documentID,
			DocumentType:    values2.DocumentTypeIdentification,
			DocumentSubType: string(values2.DocumentSubTypeRg),
			DocumentFields:  entity2.DocumentFields{Number: "888", Name: "Lola"},
		}},
		DocumentFiles: []entity2.DocumentFile{
			{DocumentFileID: &frontFileID, FileID: uuid.New(), FileSide: values2.FileSideFront, DocumentID: documentID},
			{DocumentFileID: &backFileID, FileID: uuid.New(), FileSide: values2.FileSideBack, DocumentID: documentID},
		},
	}
	analyzer = NewDoaAnalyzer(fileService, doaService, person, nil, nil)

	lastResult := &entity2.DOAResult{
		FileIDs: []uuid.UUID{
			person.DocumentFiles[0].FileID,
			person.DocumentFiles[1].FileID,
		},
		Scores: []entity2.Score{
			{FileID: person.DocumentFiles[0].FileID},
			{FileID: person.DocumentFiles[1].FileID}}}
	data := &contracts.DocumentMetadata{}

	doaService.On("FindLastResult", &person.ProfileID, &documentID).Return(lastResult, nil)

	fileService.On("GetUrl", person.DocumentFiles[0].FileID).Return("/url1", nil)

	fileService.On("GetUrl", person.DocumentFiles[1].FileID).Return("/url2", nil)

	doaService.On("CreateDocMetadata", lastResult, string(person.Documents[0].DocumentType), string(person.Documents[0].DocumentSubType), person.Documents[0].DocumentID.String(), profileID.String()).Return(data)

	_, err := analyzer.Analyze()

	assert.Nil(t, err)
}

func TestGetName(t *testing.T) {

	received := analyzer.Name()

	assert.Equal(t, values2.RuleSetDOA, received)
}

// ------------------------------------------------------------------------------------------
func TestFilesChanged(t *testing.T) {
	frontFileID := uuid.New()
	backFileID := uuid.New()
	lastResult := &entity2.DOAResult{
		Scores: entity2.Scores{
			{FileID: frontFileID},
			{FileID: uuid.New()},
		},
	}

	received := filesChanged(&frontFileID, &backFileID, lastResult)

	assert.Equal(t, true, received)

}

func TestFilesChanged_noLastResult(t *testing.T) {
	frontFileID := uuid.New()
	backFileID := uuid.New()
	var lastResult *entity2.DOAResult = nil

	received := filesChanged(&frontFileID, &backFileID, lastResult)

	assert.Equal(t, true, received)

}

func TestFilesChanged_notChanged(t *testing.T) {
	frontFileID := uuid.New()
	backFileID := uuid.New()
	lastResult := &entity2.DOAResult{
		FileIDs: []uuid.UUID{
			frontFileID,
			backFileID,
		},
	}

	received := filesChanged(&frontFileID, &backFileID, lastResult)

	assert.Equal(t, false, received)

}

func TestExistsInResult(t *testing.T) {
	frontFileID := uuid.New()
	backFileID := uuid.New()
	lastResult := &entity2.DOAResult{
		FileIDs: []uuid.UUID{
			frontFileID,
			backFileID,
		},
	}

	received := existsInResult(frontFileID, backFileID, lastResult)

	assert.True(t, received)
}

func TestExistsInResult_notEnoughScores(t *testing.T) {
	frontFileID := uuid.New()
	backFileID := uuid.New()
	lastResult := &entity2.DOAResult{
		FileIDs: []uuid.UUID{
			frontFileID,
		},
	}

	received := existsInResult(frontFileID, backFileID, lastResult)

	assert.False(t, received)
}

func TestExistsInResult_onlyContainsOne(t *testing.T) {
	frontFileID := uuid.New()
	backFileID := uuid.New()
	lastResult := &entity2.DOAResult{
		FileIDs: []uuid.UUID{
			frontFileID,
			uuid.New(),
		},
	}

	received := existsInResult(frontFileID, backFileID, lastResult)

	assert.False(t, received)
}

func TestValidateScore_noApprovedScore(t *testing.T) {
	scores := entity2.Scores{}

	received, pending := validateScore(scores, nil, nil)

	assert.Equal(t, values2.ResultStatusAnalysing, received)

	assert.True(t, pending)
}

func TestValidateScore_approved(t *testing.T) {
	scores := entity2.Scores{{Total: 0.8}, {Total: 0.6}}
	approvedScore := 0.69

	received, pending := validateScore(scores, &approvedScore, nil)

	assert.Equal(t, values2.ResultStatusApproved, received)
	assert.False(t, pending)
}

func TestValidateScore_rejected(t *testing.T) {
	scores := entity2.Scores{{Total: 0.4}, {Total: 0.2}}
	approvedScore := 0.69
	rejectedScore := 0.31

	received, pending := validateScore(scores, &approvedScore, &rejectedScore)

	assert.Equal(t, values2.ResultStatusRejected, received)
	assert.False(t, pending)
}

func TestValidateScore_noRejectedScore(t *testing.T) {
	scores := entity2.Scores{{Total: 0.4}, {Total: 0.2}}
	approvedScore := 0.69

	received, pending := validateScore(scores, &approvedScore, nil)

	assert.Equal(t, values2.ResultStatusAnalysing, received)
	assert.True(t, pending)
}

func TestContains(t *testing.T) {
	value := uuid.MustParse("7d2d0cc3-f173-4a78-86c8-5b9de2fbf8d2")
	ids := []uuid.UUID{uuid.MustParse("7d2d0cc3-f173-4a78-86c8-5b9de2fbf8d2")}

	received := contains(ids, value)

	assert.True(t, received)
}

func TestContains_notContains(t *testing.T) {
	value := uuid.MustParse("7d2d0cc3-f173-4a78-86c8-5b9de2fbf8d1")
	ids := []uuid.UUID{uuid.MustParse("7d2d0cc3-f173-4a78-86c8-5b9de2fbf8d2")}

	received := contains(ids, value)

	assert.False(t, received)
}
