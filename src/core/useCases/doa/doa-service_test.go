package doa

import (
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	mocks "bitbucket.org/bexstech/temis-compliance/src/core/_mocks"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/doa/contracts"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var (
	validate       = validator.New()
	repository     *mocks.DOAResultRepository
	eventPublisher *mocks.ComplianceCommandPublisher
	doaAdapter     *mocks.DOAAdapter
	service        _interfaces.DOAService
)

func TestMain(m *testing.M) {
	repository = &mocks.DOAResultRepository{}
	eventPublisher = &mocks.ComplianceCommandPublisher{}
	doaAdapter = &mocks.DOAAdapter{}
	service = NewDOAService(validate, repository, eventPublisher, doaAdapter)
	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	id := uuid.New()
	result := &entity2.DOAResult{}

	repository.On("Get", &id).Return(result, nil)

	received, err := service.Get(&id)

	if !reflect.DeepEqual(result, received) {
		t.Errorf("\nExpected doa result: %v \nGot: %v\n", result, received)
	}

	if err != nil {
		t.Errorf("\nExpected error nil\n")
	}
}

func TestSave(t *testing.T) {
	doaResult := &entity2.DOAResult{ID: uuid.New()}

	saved := &entity2.DOAResult{ID: uuid.New()}
	repository.On("Save", doaResult).Return(saved, nil)

	eventPublisher.On("SendCommand", saved.ID, "", values2.EventTypeDoaResultChanged).Return(nil)

	received, err := service.Save(doaResult)

	if !reflect.DeepEqual(saved, received) {
		t.Errorf("\nExpected doa result: %v \nGot: %v\n", saved, received)
	}

	if err != nil {
		t.Errorf("\nExpected error nil\n")
	}
}

func TestEnrichWithScores_success(t *testing.T) {
	id := uuid.New()
	scores := entity2.Scores{
		{FileID: uuid.New()},
	}

	doaResult := &entity2.DOAResult{}

	repository.On("Get", &id).Return(doaResult, nil)

	doaResultUpdated := &entity2.DOAResult{}
	doaResultUpdated.Scores = scores
	doaResultUpdated.Status = values2.DOAStatusDone

	saved := &entity2.DOAResult{}
	var state *entity2.State = nil

	repository.On("Save", doaResult).Return(saved, nil)

	eventPublisher.On("SendCommand", saved.ID, &saved.EntityID, values2.EngineNameProfile).Return(state, nil)

	received, err := service.EnrichWithScores(&id, scores)

	if !reflect.DeepEqual(saved, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", saved, received)
	}

	if err != nil {
		t.Errorf("\nExpected error nil\n")
	}
}

func TestEnrichWithScores_errGettingDoaResult(t *testing.T) {
	id := uuid.New()
	scores := entity2.Scores{
		{FileID: uuid.New()},
	}

	doaResult := &entity2.DOAResult{}

	errGettingDoaResult := errors.New("error getting doa result")
	repository.On("Get", &id).Return(doaResult, errGettingDoaResult)

	received, err := service.EnrichWithScores(&id, scores)

	if received != nil {
		t.Errorf("\nExpected doa result: %v \n", nil)
	}

	if !reflect.DeepEqual(errGettingDoaResult, err) {
		t.Errorf("\nExpected: %v \nGot: %v\n", errGettingDoaResult, err)
	}
}

func TestEnrichWithScores_noResult(t *testing.T) {
	id := uuid.New()
	scores := entity2.Scores{
		{FileID: uuid.New()},
	}
	var doaResult *entity2.DOAResult = nil
	repository.On("Get", &id).Return(doaResult, nil)

	received, err := service.EnrichWithScores(&id, scores)

	if received != nil {
		t.Errorf("\nExpected doa result: %v \n", nil)
	}

	if err != nil {
		t.Errorf("\nExpected error: %v \n", nil)
	}
}

func TestEnrichWithScores_errSavingDoaResult(t *testing.T) {
	id := uuid.New()
	scores := entity2.Scores{
		{FileID: uuid.New()},
	}

	doaResult := &entity2.DOAResult{}

	repository.On("Get", &id).Return(doaResult, nil)

	doaResultUpdated := &entity2.DOAResult{}
	doaResultUpdated.Scores = scores
	doaResultUpdated.Status = values2.DOAStatusDone

	saved := &entity2.DOAResult{}
	errSavingDoaResult := errors.New("error saving doa result")
	repository.On("Save", doaResultUpdated).Return(saved, errSavingDoaResult)

	received, err := service.EnrichWithScores(&id, scores)

	if received != nil {
		t.Errorf("\nExpected doa result: %v \n", nil)
	}

	if !reflect.DeepEqual(errSavingDoaResult, err) {
		t.Errorf("\nExpected: %v \nGot: %v\n", errSavingDoaResult, err)
	}
}

func TestEnrichWithScores_errSendingEvent(t *testing.T) {
	id := uuid.New()
	scores := entity2.Scores{
		{FileID: uuid.New()},
	}

	doaResult := &entity2.DOAResult{}

	repository.On("Get", &id).Return(doaResult, nil)

	doaResultUpdated := &entity2.DOAResult{}
	doaResultUpdated.Scores = scores
	doaResultUpdated.Status = values2.DOAStatusDone

	saved := &entity2.DOAResult{ID: uuid.New()}
	var state *entity2.State
	repository.On("Save", doaResultUpdated).Return(saved, nil)

	errSendingEvent := errors.New("error sending event")
	eventPublisher.On("SendCommand", saved.ID, &saved.EntityID, values2.EngineNameProfile).Return(state, errSendingEvent)
	received, err := service.EnrichWithScores(&id, scores)

	if received != nil {
		t.Errorf("\nExpected doa result: %v \n", nil)
	}

	if !reflect.DeepEqual(errSendingEvent, err) {
		t.Errorf("\nExpected: %v \nGot: %v\n", errSendingEvent, err)
	}
}

func TestFindLastResult_success(t *testing.T) {
	entityID := uuid.New()
	documentID := uuid.New()

	lastResult := &entity2.DOAResult{ID: uuid.New()}

	repository.On("FindLastByEntityIdAndDocumentId", &entityID, &documentID).Return(lastResult, nil)

	expected := lastResult
	received, err := service.FindLastResult(&entityID, &documentID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v \nGot: %v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected error: %v \n", nil)
	}
}

func TestFindLastResult_errFinding(t *testing.T) {
	entityID := uuid.New()
	documentID := uuid.New()

	lastResult := &entity2.DOAResult{ID: uuid.New()}

	errFinding := errors.New("err finding")
	repository.On("FindLastByEntityIdAndDocumentId", &entityID, &documentID).Return(lastResult, errFinding)

	received, err := service.FindLastResult(&entityID, &documentID)

	if !reflect.DeepEqual(errFinding, err) {
		t.Errorf("\nExpected: %v \nGot: %v\n", errFinding, err)
	}

	if received != nil {
		t.Errorf("\nExpected result nil\n")
	}
}

func TestRequestExtraction(t *testing.T) {
	frontFile := &entity2.DocumentFile{FileID: uuid.New()}
	frontFileURI := "frontFileURI"
	backFile := &entity2.DocumentFile{FileID: uuid.New()}
	backFileURI := "backFileURL"
	time := time.Date(2020, 10, 10, 10, 10, 10, 10, time.UTC)
	doc := &entity2.Document{
		DocumentSubType: string(values2.DocumentSubTypeRg),
		DocumentFields:  entity2.DocumentFields{Number: "12345", IssueDate: time.String(), Name: "Maria"},
	}
	profileID := uuid.New()
	extraction := &entity2.DOAExtraction{RequestID: uuid.New()}

	doaAdapter.On("RequestExtraction", frontFile, frontFileURI, backFile, backFileURI, doc, profileID).Return(extraction, nil)

	expected := extraction
	received, err := service.RequestExtraction(frontFile, frontFileURI, backFile, backFileURI, doc, profileID)

	if !reflect.DeepEqual(expected, received) {
		t.Errorf("\nExpected: %v\nGot:%v\n", expected, received)
	}

	if err != nil {
		t.Errorf("\nExpected error nil\n")
	}
}

func TestCreateDocMetadata(t *testing.T) {
	doaResult := &entity2.DOAResult{
		Scores: []entity2.Score{
			{
				FileID:          uuid.New(),
				ForDocumentSide: entity2.ScoredDocumentSide{Given: "FRONT"},
				Total:           0.99,
			},
			{
				FileID:          uuid.New(),
				ForDocumentSide: entity2.ScoredDocumentSide{Given: "BACK"},
				Total:           0.90,
			},
		},
	}
	documentType := "Identification"
	docSubType := "RG"
	documentID := "111"
	profileID := "222"

	expectedMetadata := &contracts.DocumentMetadata{
		Type:       documentType,
		SubType:    docSubType,
		DocumentID: documentID,
		ProfileID:  profileID,
		Files: []contracts.File{
			{
				Side:       "FRONT",
				FileID:     doaResult.Scores[0].FileID.String(),
				TotalScore: 0.99,
			},
			{
				Side:       "BACK",
				FileID:     doaResult.Scores[1].FileID.String(),
				TotalScore: 0.90,
			},
		},
	}

	receivedMetadata := service.CreateDocMetadata(doaResult, documentType, docSubType,
		documentID, profileID)

	if !reflect.DeepEqual(expectedMetadata, receivedMetadata) {
		t.Errorf("\nExpected: %v\nGot:%v\n", expectedMetadata, receivedMetadata)
	}
}
