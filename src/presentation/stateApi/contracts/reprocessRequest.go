package contracts

type ReprocessRequest struct {
	Ids        []string `json:"ids"`
	EngineName string   `json:"engine_name"`
}
