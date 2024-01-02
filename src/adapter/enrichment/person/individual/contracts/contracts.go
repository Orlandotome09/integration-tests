package contracts

type IndividualResponse struct {
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
	Situation int    `json:"situation"`
}
