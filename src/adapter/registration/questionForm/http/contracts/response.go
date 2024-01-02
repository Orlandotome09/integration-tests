package contracts

type QuestionFormResponse struct {
	QuestionFormID string `json:"question_form_id"`
	EntityID       string `json:"entity_id"`
	Code           string `json:"code,omitempty"`
	Answer         string `json:"answer,omitempty"`
	Comments       string `json:"comments,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}
