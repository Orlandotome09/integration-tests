package person

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"github.com/google/uuid"
)

type doaAnalyzer struct {
	fileService   interfaces.FileAdapter
	doaService    interfaces.DOAService
	person        entity2.Person
	approvedScore *float64
	rejectedScore *float64
}

func NewDoaAnalyzer(
	fileService interfaces.FileAdapter,
	doaService interfaces.DOAService,
	person entity2.Person,
	approvedScore *float64,
	rejectedScore *float64,
) entity2.Rule {
	return &doaAnalyzer{
		fileService:   fileService,
		doaService:    doaService,
		person:        person,
		approvedScore: approvedScore,
		rejectedScore: rejectedScore,
	}
}

func (ref *doaAnalyzer) Analyze() ([]entity2.RuleResultV2, error) {
	documentNotFound := entity2.NewRuleResultV2(values2.RuleSetDOA, values2.RuleNameDocumentNotFound)
	profileID := ref.person.EntityID

	documents := make([]entity2.Document, 0)

	foundIdentificationDocument := false

	for _, document := range ref.person.Documents {
		if document.DocumentType == values2.DocumentTypeIdentification {
			foundIdentificationDocument = true
			documents = append(documents, document)
		}
	}

	if !foundIdentificationDocument {
		metadata, _ := json.Marshal("Identification Document Not Found")
		documentNotFound.SetResult(values2.ResultStatusIncomplete).SetMetadata(metadata).
			AddProblem(values2.ProblemCodeIdentificationDocumentNotFound, "")

		return []entity2.RuleResultV2{*documentNotFound}, nil
	}
	documentNotFound.SetResult(values2.ResultStatusApproved)

	rulesResult := []entity2.RuleResultV2{*documentNotFound}
	for _, document := range documents {
		ruleResults, err := ref.analyzeDocument(profileID, document)
		if err != nil {
			return nil, err
		}
		rulesResult = append(rulesResult, ruleResults...)
	}
	return rulesResult, nil
}

func (ref *doaAnalyzer) analyzeDocument(profileID uuid.UUID, doc entity2.Document) ([]entity2.RuleResultV2, error) {
	doaFileNotFound := entity2.NewRuleResultV2(values2.RuleSetDOA, values2.RuleNameDOAFileNotfound)

	var frontFile entity2.DocumentFile
	var backFile entity2.DocumentFile

	for _, file := range ref.person.DocumentFiles {

		if file.DocumentID == doc.DocumentID {
			if file.FileSide == values2.FileSideFront {
				frontFile = file
			}
			if file.FileSide == values2.FileSideBack {
				backFile = file
			}
		}
	}

	if frontFile.DocumentFileID == nil {
		metadata, _ := json.Marshal("FRONT FILE Not Found")
		doaFileNotFound.SetResult(values2.ResultStatusIncomplete).SetMetadata(metadata).
			AddProblem(values2.ProblemCodeFrontFileNotFound, "")
		return []entity2.RuleResultV2{*doaFileNotFound}, nil
	}

	if backFile.DocumentFileID == nil {
		metadata, _ := json.Marshal("BACK FILE Not Found")
		doaFileNotFound.SetResult(values2.ResultStatusIncomplete).SetMetadata(metadata).
			AddProblem(values2.ProblemCodeBackFileNotFound, "")
		return []entity2.RuleResultV2{*doaFileNotFound}, nil
	}

	doaFileNotFound.SetResult(values2.ResultStatusApproved)
	lastResult, err := ref.doaService.FindLastResult(&profileID, &doc.DocumentID)
	if err != nil {
		return nil, err
	}

	frontFileURI, err := ref.fileService.GetUrl(frontFile.FileID)
	if err != nil {
		return nil, err
	}

	backFileURI, err := ref.fileService.GetUrl(backFile.FileID)
	if err != nil {
		return nil, err
	}

	changed := filesChanged(&frontFile.FileID, &backFile.FileID, lastResult)

	doaValidation := entity2.NewRuleResultV2(values2.RuleSetDOA, values2.RuleNameDOAValidation)
	if changed {

		resp, err := ref.doaService.RequestExtraction(&frontFile, frontFileURI, &backFile, backFileURI,
			&doc, profileID)
		if err != nil {
			return nil, err
		}

		newDoaResult := &entity2.DOAResult{ID: resp.RequestID,
			EntityID:   profileID,
			DocumentID: doc.DocumentID,
			FileIDs:    []uuid.UUID{frontFile.FileID, backFile.FileID},
			Status:     values2.DOAStatusValidating,
		}

		_, err = ref.doaService.Save(newDoaResult)
		if err != nil {
			return nil, err
		}

		metadata, _ := json.Marshal("Awaiting for DOA calculate score")

		doaValidation.SetResult(values2.ResultStatusAnalysing).SetMetadata(metadata)

		return []entity2.RuleResultV2{*doaFileNotFound, *doaValidation}, nil
	}

	data := ref.doaService.CreateDocMetadata(lastResult, string(doc.DocumentType), string(doc.DocumentSubType), doc.DocumentID.String(), profileID.String())

	metadata, _ := json.Marshal(data)

	status, pending := validateScore(lastResult.Scores, ref.approvedScore, ref.rejectedScore)

	doaValidation.SetResult(status).SetPending(pending).SetMetadata(metadata)
	return []entity2.RuleResultV2{*doaFileNotFound, *doaValidation}, nil
}

func (ref *doaAnalyzer) Name() values2.RuleSet {
	return values2.RuleSetDOA
}

// ---------------------------------------------------------------------------------------------------------
func filesChanged(frontFileID *uuid.UUID, backFileID *uuid.UUID,
	lastResult *entity2.DOAResult) bool {

	if lastResult == nil {
		return true
	}

	exists := existsInResult(*frontFileID, *backFileID, lastResult)

	changed := !exists

	return changed
}

func existsInResult(frontFileID uuid.UUID, backFileID uuid.UUID, lastResult *entity2.DOAResult) bool {

	if len(lastResult.FileIDs) != 2 {
		return false
	}

	containsFrontFileID := contains(lastResult.FileIDs, frontFileID)
	containsBackFileID := contains(lastResult.FileIDs, backFileID)

	if containsFrontFileID && containsBackFileID {
		return true
	}
	return false
}

func validateScore(scores entity2.Scores, approvedScore *float64,
	rejectedScore *float64) (values2.Result, bool) {

	approval := approvedScore
	reject := rejectedScore

	if approval == nil {
		return values2.ResultStatusAnalysing, true
	}

	totalScore := 0.0

	for _, score := range scores {
		totalScore = totalScore + score.Total
	}

	meanScore := totalScore / float64(len(scores))

	if meanScore > *approval {
		return values2.ResultStatusApproved, false
	}

	if reject != nil {
		if meanScore < *reject {
			return values2.ResultStatusRejected, false
		}
	}

	return values2.ResultStatusAnalysing, true

}

func contains(ids []uuid.UUID, value uuid.UUID) bool {
	for _, id := range ids {
		if value == id {
			return true
		}
	}

	return false
}
