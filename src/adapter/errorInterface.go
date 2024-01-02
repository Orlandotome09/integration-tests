package adapter

import "github.com/pkg/errors"

type ErrorResponse struct {
	Error string `json:"error"`
}

func (errorResponse *ErrorResponse) ToError() error {

	return errors.Errorf("Error: %s", errorResponse.Error)
}

type ErrorInterface interface {
	ToError() error
}