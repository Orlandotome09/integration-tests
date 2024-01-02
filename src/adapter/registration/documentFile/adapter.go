package documentFile

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile/http"
	"bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile/http/contracts"
	translator "bitbucket.org/bexstech/temis-compliance/src/adapter/registration/documentFile/translator"
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"sort"
)

const (
	frontFile = "FRONT"
	backFile  = "BACK"
)

type documentFileAdapter struct {
	documentFileClient documentFileClient.DocumentFileClient
	translator         translator.DocumentFileTranslator
}

func NewDocumentFileAdapter(documentFileClient documentFileClient.DocumentFileClient,
	translator translator.DocumentFileTranslator) interfaces.DocumentFileAdapter {
	return &documentFileAdapter{
		documentFileClient: documentFileClient,
		translator:         translator,
	}
}

func (ref *documentFileAdapter) Get(documentFileID uuid.UUID) (*entity.DocumentFile, error) {
	resp, err := ref.documentFileClient.Get(documentFileID.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp == nil {
		return nil, nil
	}

	documentFile := ref.translator.Translate(*resp)

	return &documentFile, nil
}

func (ref *documentFileAdapter) FindByDocumentID(documentID uuid.UUID) ([]entity.DocumentFile, error) {
	resp, err := ref.documentFileClient.SearchByDocumentId(documentID.String())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	documentFiles := ref.translator.TranslateAll(resp)

	return documentFiles, nil
}
func (ref *documentFileAdapter) GetLastTwoFilesOfDocument(documentID uuid.UUID) (*entity.DocumentFile,
	*entity.DocumentFile, error) {
	resp, err := ref.documentFileClient.SearchByDocumentId(documentID.String())
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	frontFileResponses, backFileResponses := splitResponsesByFileSide(resp)

	sortedFront, sortedBack := sortResponsesDescending(frontFileResponses), sortResponsesDescending(backFileResponses)

	frontDocumentFiles, backDocumentFiles := ref.translator.TranslateAll(sortedFront), ref.translator.TranslateAll(sortedBack)

	frontFile, backFile := findLast(frontDocumentFiles), findLast(backDocumentFiles)

	return frontFile, backFile, nil
}

// ------------------------------------------------------------------------
func splitResponsesByFileSide(responses []contracts.DocumentFileResponse) ([]contracts.DocumentFileResponse,
	[]contracts.DocumentFileResponse) {
	frontFileResponses := []contracts.DocumentFileResponse{}
	backFileResponses := []contracts.DocumentFileResponse{}

	for _, elem := range responses {
		if elem.FileSide == frontFile {
			frontFileResponses = append(frontFileResponses, elem)
		}
		if elem.FileSide == backFile {
			backFileResponses = append(backFileResponses, elem)
		}
	}

	return frontFileResponses, backFileResponses

}

func sortResponsesDescending(responses []contracts.DocumentFileResponse) []contracts.DocumentFileResponse {
	sort.Slice(responses, func(i, j int) bool {
		return responses[i].CreatedAt.After(responses[j].CreatedAt)
	})

	return responses
}

func findLast(sortedDocumentFiles []entity.DocumentFile) *entity.DocumentFile {

	if len(sortedDocumentFiles) == 0 {
		return nil
	}

	return &sortedDocumentFiles[0]
}
