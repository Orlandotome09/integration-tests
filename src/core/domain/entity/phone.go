package entity

type Phone struct {
	Type        string `json:"type"`
	CountryCode string `json:"country_code"`
	AreaCode    string `json:"area_code"`
	Number      string `json:"number"`
}
