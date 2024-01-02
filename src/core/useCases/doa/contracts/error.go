package contracts

import "fmt"

type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func (ref *ErrorResponse) Error() string {
	return fmt.Sprintf("error message: %v , detail: %v ", ref.Message, ref.Details)
}
