package entity

import (
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
)

type DocumentRequired struct {
	DocumentType      values.DocumentType `json:"document_type"`
	DocumentSubType   string              `json:"document_sub_type"`
	FileRequired      bool                `json:"file_required"`
	PendingOnApproval bool                `json:"pending_on_approval"`
	Conditions        []Condition         `json:"conditions"`
}

func (documentRequired DocumentRequired) HasSubtype() bool {
	return documentRequired.DocumentSubType != ""
}

type DocumentsRequired []DocumentRequired

func (documentsRequired DocumentsRequired) HaveAnyPendingOnApproval() bool {
	for _, documentRequired := range documentsRequired {
		if documentRequired.PendingOnApproval {
			return true
		}
	}
	return false
}

func (documentsRequired DocumentsRequired) FilterByConditionsSatisfied(person Person) DocumentsRequired {
	filtered := DocumentsRequired{}
	for _, documentRequired := range documentsRequired {
		if person.SatisfiesAllConditions(documentRequired.Conditions) {
			filtered = append(filtered, documentRequired)
		}
	}

	return filtered
}
