package message

type Individual struct {
	Pep    bool    `json:"pep" binding:"required"`
	Income float64 `json:"income" binding:"required"`
	Assets float64 `json:"assets" binding:"required"`
}
