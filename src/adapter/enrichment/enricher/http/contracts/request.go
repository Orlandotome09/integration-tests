package enricherContracts

type EnricherRequest struct {
	ProfileID  string `json:"profile_id"`
	PersonType string `json:"person_type"`
	OfferType  string `json:"offer_type"`
	PartnerID  string `json:"partner_id"`
	RoleType   string `json:"role_type"`
}
