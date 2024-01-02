package fileHttpClient

type FileResponse struct {
	FileID string `json:"file_id,omitempty"`
	Url    string `json:"url,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
