package message

type LimitMessage struct {
	EventType      string   `json:"event_type,omitempty""`
	ProfileID      string   `json:"profile_id" binding:"required"`
	DocumentNumber string   `json:"document_number,omitempty"`
	PartnerID      string   `json:"partner_id,omitempty"`
	OfferType      string   `json:"offer_type" binding:"required"`
	PersonType     string   `json:"person_type" binding:"required"`
	RoleType       string   `json:"role_type" binding:"required"`
	ApprovedRules  []string `json:"approved_rules" binding:"required"`
	Documents      []string `json:"documents,omitempty"`
}

var EventTypeProfileApproved string = "PROFILE_APPROVED"
