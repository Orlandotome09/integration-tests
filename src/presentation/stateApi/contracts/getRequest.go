package contracts

type GetRequest struct {
	ProfileID   string `uri:"profile_id" binding:"uuid"`
	OnlyPending bool   `form:"only_pending"`
}

type GetComplianceRequest struct {
	EntityID string `uri:"entity_id" binding:"uuid"`
}

type GetComplianceCheckRequest struct {
	EntityID   string `uri:"entity_id" binding:"uuid"`
	EngineName string `json:"entity_type"`
	EventType  string `json:"event_type,omitempty"`
}
