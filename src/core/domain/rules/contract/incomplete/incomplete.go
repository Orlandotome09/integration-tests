package incompleteContractRule

import (
	interfaces "bitbucket.org/bexstech/temis-compliance/src/core/_interfaces"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/rules/contract"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"encoding/json"
	"fmt"
)

type incompleteContractRule struct {
	contractRule.ContractRule
	documentAdapter     interfaces.DocumentAdapter
	documentFileAdapter interfaces.DocumentFileAdapter
}

func New(contract entity.Contract,
	documentAdapter interfaces.DocumentAdapter,
	documentFileAdapter interfaces.DocumentFileAdapter) entity.Rule {
	return &incompleteContractRule{
		ContractRule:        contractRule.ContractRule{Contract: contract},
		documentAdapter:     documentAdapter,
		documentFileAdapter: documentFileAdapter,
	}
}

func (ref *incompleteContractRule) Analyze() ([]entity.RuleResultV2, error) {
	documentNotFound := entity.NewRuleResultV2(values.RuleSetIncompleteContract, values.RuleNameDocumentNotFound)
	fileNotFound := entity.NewRuleResultV2(values.RuleSetIncompleteContract, values.RuleNameFileNotFound)

	if ref.Contract.DocumentID == nil {
		metadata, _ := json.Marshal(fmt.Sprintf("Invoice is required"))
		documentNotFound.
			SetResult(values.ResultStatusIncomplete).
			SetMetadata(metadata).
			AddProblem(values.ProblemCodeInvoiceIsRequired, "")
		return []entity.RuleResultV2{*documentNotFound, *fileNotFound}, nil
	}

	documentID := *ref.Contract.DocumentID
	document, err := ref.documentAdapter.GetByID(documentID.String())
	if err != nil {
		return nil, err
	}
	if document == nil {
		metadata, _ := json.Marshal(fmt.Sprintf("Invoice Document: %v not found", documentID))
		documentNotFound.
			SetResult(values.ResultStatusIncomplete).
			SetMetadata(metadata).
			AddProblem(values.ProblemCodeInvoiceDocumentNotFound, documentID.String())
		return []entity.RuleResultV2{*documentNotFound, *fileNotFound}, nil
	}

	//Reject document if it is associated to another profile
	if *ref.Contract.ProfileID != document.EntityID {
		metadata, _ := json.Marshal(fmt.Sprintf("Invoice Document is associated to another profile"))
		documentNotFound.
			SetResult(values.ResultStatusRejected).
			SetMetadata(metadata).
			AddProblem(values.ProblemCodeInvoiceAssociatedToAnotherProfile, documentID.String())
		return []entity.RuleResultV2{*documentNotFound, *fileNotFound}, nil
	}

	documentNotFound.SetResult(values.ResultStatusApproved)

	documentFiles, err := ref.documentFileAdapter.FindByDocumentID(documentID)
	if err != nil {
		return nil, err
	}
	if len(documentFiles) == 0 {
		metadata, _ := json.Marshal(fmt.Sprintf("Invoice File not found for Document: %v",
			documentID))

		fileNotFound.
			SetResult(values.ResultStatusIncomplete).
			SetMetadata(metadata).
			AddProblem(values.ProblemCodeInvoiceFileNotFound, documentID.String())
		return []entity.RuleResultV2{*documentNotFound, *fileNotFound}, nil
	}
	fileNotFound.SetResult(values.ResultStatusApproved).SetPending(true)

	return []entity.RuleResultV2{*documentNotFound, *fileNotFound}, nil
}

func (ref *incompleteContractRule) Name() values.RuleSet {
	return values.RuleSetIncompleteContract
}
