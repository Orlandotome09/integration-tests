package contracts

type PartnerResponse struct {
	PartnerID      string `json:"partner_id"`
	DocumentNumber string `json:"document_number" `
	Name           string `json:"name"`
	Status         string `json:"status"`
	LogoImageUrl   string `json:"logo_image_url"`
	Config         Config `json:"config"`
}

type Config struct {
	CustomerSegregationType string `json:"customer_segregation_type"`
	UseCallbackV2           bool   `json:"use_callback_v2,omitempty"`
}
