package authentication

import (
	"errors"
	"net/http"
)

type Provider interface {
	Validate(r *http.Request) error
}

var (
	provider Provider

	ErrInvalidToken    = errors.New("Token is not valid")
	ErrExpiredToken    = errors.New("Token is expired")
	errSettingsMissing = errors.New("Missing one of required env variables: AUTH0_JWKS_URI, AUTH0_AUDIENCE and AUTH0_ISSUER")
)

func Setup(p Provider) {
	provider = p
}

func Validate(r *http.Request) error {
	err := provider.Validate(r)
	return err
}
