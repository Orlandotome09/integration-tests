package limitMessageTranslator

import (
	"bitbucket.org/bexstech/temis-compliance/src/adapter/event/publisher/limit/message"
	entity2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/entity"
	values2 "bitbucket.org/bexstech/temis-compliance/src/core/domain/values"
	"github.com/google/uuid"
)

type LimitMessageTranslator interface {
	Translate(profile entity2.Profile, state entity2.State) *message.LimitMessage
}

type limitMessageTranslator struct{}

func New() LimitMessageTranslator {
	return &limitMessageTranslator{}
}

func (ref *limitMessageTranslator) Translate(profile entity2.Profile, state entity2.State) *message.LimitMessage {
	return &message.LimitMessage{
		EventType:      message.EventTypeProfileApproved,
		ProfileID:      translateProfileID(profile.ProfileID),
		DocumentNumber: profile.DocumentNumber,
		PartnerID:      profile.PartnerID,
		OfferType:      profile.OfferType,
		PersonType:     translatePersonType(profile.Person.PersonType),
		RoleType:       translateRoleType(profile.RoleType),
		ApprovedRules:  translateApprovedRule(state.ValidationStepsResults.FindApprovedRules()),
		Documents:      translateDocuments(profile),
	}
}

func translateProfileID(profileID *uuid.UUID) string {
	if profileID == nil {
		return ""
	}
	return profileID.String()
}

func translatePersonType(personType values2.PersonType) string {
	return string(personType)
}

func translateRoleType(roleType values2.RoleType) string {
	return string(roleType)
}

func translateDocuments(profile entity2.Profile) []string {
	documents := make([]string, len(profile.Documents))

	for i, document := range profile.Documents {
		documents[i] = document.DocumentType
	}

	return documents
}

func translateApprovedRule(approvedRules []values2.RuleName) []string {
	rules := make([]string, len(approvedRules))

	for i, approvedRule := range approvedRules {
		rules[i] = approvedRule.ToString()
	}

	return rules
}
