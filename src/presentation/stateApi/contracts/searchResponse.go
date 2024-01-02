package contracts

import (
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	"bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"bitbucket.org/bexstech/temis-compliance/src/presentation"
	"github.com/google/uuid"
	"time"
)

type SearchResponse struct {
	ProfileID      uuid.UUID      `json:"profile_id"`
	PartnerID      string         `json:"partner_id"`
	ParentID       *uuid.UUID     `json:"parent_id,omitempty"`
	DocumentNumber string         `json:"document_number"`
	Name           string         `json:"name"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Rules          []RuleResponse `json:"rules"`
}

func (ref SearchResponse) FromDomain(profileState entity2.ProfileState, onlyPending bool) SearchResponse {
	ref.ProfileID = profileState.ProfileID
	ref.PartnerID = profileState.PartnerID
	ref.ParentID = profileState.ParentID
	ref.DocumentNumber = profileState.DocumentNumber
	ref.Name = profileState.Name
	ref.CreatedAt = profileState.State.CreatedAt
	ref.UpdatedAt = profileState.State.UpdatedAt
	ref.Rules = RuleResponse{}.FromDomainValidationSteps(profileState.State.ValidationStepsResults, onlyPending)
	return ref
}

type RuleResponse struct {
	Name values.RuleName `json:"name"`
}

func (ref RuleResponse) FromDomainValidationSteps(validationStepResults []entity2.ValidationStepResult, onlyPending bool) []RuleResponse {

	rulesResponse := make([]RuleResponse, 0)

	for _, step := range validationStepResults {
		for _, ruleResult := range step.RuleResults {
			if !onlyPending || ruleResult.Pending {
				rulesResponse = append(rulesResponse, RuleResponse{
					Name: ruleResult.RuleName,
				})
			}
		}
	}

	return rulesResponse
}

type SearchPaginatedResponse struct {
	Profiles []SearchResponse `json:"profiles"`
	presentation.PagingResponse
}
