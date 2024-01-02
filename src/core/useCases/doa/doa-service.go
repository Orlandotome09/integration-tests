package doa

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/core/useCases/doa/contracts"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type doaService struct {
	validate       *validator.Validate
	repository     interfaces.DOAResultRepository
	eventPublisher interfaces.ComplianceCommandPublisher
	adapter        interfaces.DOAAdapter
}

func NewDOAService(validate *validator.Validate, repository interfaces.DOAResultRepository,
	eventPublisher interfaces.ComplianceCommandPublisher, adapter interfaces.DOAAdapter,
) interfaces.DOAService {
	return &doaService{
		validate:       validate,
		repository:     repository,
		eventPublisher: eventPublisher,
		adapter:        adapter,
	}
}

func (ref *doaService) Get(id *uuid.UUID) (*entity2.DOAResult, error) {
	return ref.repository.Get(id)
}

func (ref *doaService) Save(doaResult *entity2.DOAResult) (*entity2.DOAResult, error) {
	return ref.repository.Save(doaResult)
}

func (ref *doaService) EnrichWithScores(id *uuid.UUID, scores entity2.Scores) (*entity2.DOAResult, error) {
	doaResult, err := ref.repository.Get(id)
	if err != nil {
		return nil, err
	}

	if doaResult == nil {
		return nil, nil
	}

	doaResult.Scores = scores

	doaResult.Status = values2.DOAStatusDone

	saved, err := ref.repository.Save(doaResult)
	if err != nil {
		return nil, err
	}

	_, err = ref.eventPublisher.SendCommand(saved.ID, &saved.EntityID, values2.EngineNameProfile)
	if err != nil {
		return nil, err
	}

	return saved, nil
}

func (ref *doaService) FindLastResult(entityID *uuid.UUID, documentID *uuid.UUID) (*entity2.DOAResult, error) {
	lastResult, err := ref.repository.FindLastByEntityIdAndDocumentId(entityID, documentID)
	if err != nil {
		return nil, err
	}
	return lastResult, nil
}

func (ref *doaService) RequestExtraction(frontFile *entity2.DocumentFile, frontFileURI string,
	backFile *entity2.DocumentFile, backFileURI string, doc *entity2.Document,
	profileID uuid.UUID,
) (*entity2.DOAExtraction, error) {
	return ref.adapter.RequestExtraction(frontFile, frontFileURI, backFile, backFileURI, doc, profileID)
}

func (ref *doaService) CreateDocMetadata(doaResult *entity2.DOAResult, documentType string,
	docSubType string, documentID string, profileID string,
) *contracts.DocumentMetadata {
	metadata := &contracts.DocumentMetadata{
		Type:       documentType,
		SubType:    docSubType,
		DocumentID: documentID,
		ProfileID:  profileID,
		Files:      convertScores(doaResult),
	}

	return metadata
}

func convertScores(lastResult *entity2.DOAResult) []contracts.File {
	files := []contracts.File{}

	for _, score := range lastResult.Scores {

		file := contracts.File{
			Side:       score.ForDocumentSide.Given,
			FileID:     score.FileID.String(),
			TotalScore: score.Total,
		}

		files = append(files, file)
	}

	return files
}
