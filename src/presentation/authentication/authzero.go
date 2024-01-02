package authentication

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	"github.com/auth0-community/go-auth0"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type AuthZeroProvider struct {
	jwtValidator *auth0.JWTValidator
}

type AuthZeroSettings struct {
	JwksURI  string
	Audience []string
	Issuer   string
}

var (
	log = logrus.WithField("package", "authentication")
)

func NewAuthZeroProvider(settings AuthZeroSettings) *AuthZeroProvider {
	if len(settings.JwksURI) <= 0 || len(settings.Audience) <= 0 || len(settings.Issuer) <= 0 {
		log.Panic(errSettingsMissing)
	}

	jwtClient := auth0.NewJWKClient(auth0.JWKClientOptions{URI: settings.JwksURI}, nil)
	jwtConfiguration := auth0.NewConfiguration(jwtClient, settings.Audience, settings.Issuer, jose.RS256)
	jwtValidator := auth0.NewValidator(jwtConfiguration, nil)

	return &AuthZeroProvider{jwtValidator: jwtValidator}
}

func (ref *AuthZeroProvider) Validate(r *http.Request) error {
	_, err := ref.validateRequestAndExtractToken(r)
	if err != nil {
		return err
	}

	return errors.WithStack(err)
}

func (ref *AuthZeroProvider) validateRequestAndExtractToken(r *http.Request) (*jwt.JSONWebToken, error) {
	token, err := ref.jwtValidator.ValidateRequest(r)
	if err == jwt.ErrExpired {
		return nil, ErrExpiredToken
	}

	if err != nil {
		return nil, ErrInvalidToken
	}

	return token, nil
}
