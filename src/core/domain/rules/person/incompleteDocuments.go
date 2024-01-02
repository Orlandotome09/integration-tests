package person

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type incompleteDocumentsAnalyzer struct {
	documentsRequired entity.DocumentsRequired
}

func NewIncompleteDocumentsAnalyzer(
	documentsRequired entity.DocumentsRequired) Analyzer {
	return &incompleteDocumentsAnalyzer{
		documentsRequired: documentsRequired,
	}
}

func (ref *incompleteDocumentsAnalyzer) Analyze(person entity.Person) (*entity.RuleResultV2, error) {
	documentNotFound := entity.NewRuleResultV2(values.RuleSetIncomplete, values.RuleNameDocumentNotFound)
	var result []string
	var personDocuments []entity.Document

	documentsRequired := ref.documentsRequired.FilterByConditionsSatisfied(person)

	for _, documentRequired := range documentsRequired {
		var foundDocuments entity.Documents

		documents := person.FindDocumentsEquals(documentRequired)
		if len(documents) > 0 {
			personDocuments = append(personDocuments, documents...)
			foundDocuments = append(foundDocuments, documents...)
		} else {
			docTypeOrSubType := documentRequired.DocumentType
			if documentRequired.HasSubtype() {
				docTypeOrSubType = documentRequired.DocumentSubType
			}
			result = append(result, fmt.Sprintf("Not found document: %s", docTypeOrSubType))
			documentNotFound.AddProblem(fmt.Sprintf("%s%s", "DOCUMENT_NOT_FOUND_", docTypeOrSubType), "")
			continue
		}

		if !foundDocuments.HaveRequiredFile(person.DocumentFiles, documentRequired.FileRequired) {
			fileNameTypeOrSubType := documentRequired.DocumentType
			if documentRequired.DocumentSubType != "" {
				fileNameTypeOrSubType = documentRequired.DocumentSubType
			}
			result = append(result, fmt.Sprintf("Not found files for document: %s", fileNameTypeOrSubType))
			documentNotFound.AddProblem(fmt.Sprintf("%s%s", "FILE_NOT_FOUND_", fileNameTypeOrSubType), "")
		}
	}

	if len(result) > 0 {
		metadata, err := json.Marshal(result)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		documentNotFound.SetResult(values.ResultStatusIncomplete).SetMetadata(metadata)
		return documentNotFound, nil
	}

	if documentsRequired.HaveAnyPendingOnApproval() {
		metadata, err := json.Marshal(personDocuments)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		documentNotFound.SetResult(values.ResultStatusApproved).SetPending(true).SetMetadata(metadata)
		return documentNotFound, nil
	}

	documentNotFound.SetResult(values.ResultStatusApproved).SetPending(false)
	return documentNotFound, nil
}
