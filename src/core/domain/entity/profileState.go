package entity

import (
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
)

type ProfileState struct {
	State          State
	ProfileID      uuid.UUID
	PartnerID      string
	ParentID       *uuid.UUID
	DocumentNumber string
	Name           string
}

type ProfileStateList struct {
	Count         int64
	ProfileStates []ProfileState
}

type SearchProfileStateRequest struct {
	values2.Paging
	ProfileID      *uuid.UUID
	PartnerIDs     []string
	ParentIDs      []uuid.UUID
	DocumentNumber string
	Name           string
	OfferTypes     []string
	RuleName       values2.RuleName
	ResultStatus   []values2.Result
}
