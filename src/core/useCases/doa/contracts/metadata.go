package contracts

type DocumentMetadata struct {
	Type       string `json:"type"`
	SubType    string `json:"sub_type"`
	DocumentID string `json:"document_id"`
	ProfileID  string `json:"profile_id"`
	Files      []File `json:"files"`
}

type File struct {
	Side       string  `json:"side"`
	FileID     string  `json:"file_id"`
	TotalScore float64 `json:"total_score"`
}
