package google

import (
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

func ConfigFromJSON(jsonKey []byte, scopes ...string) (*jwt.Config, error) {

	c, err := google.JWTConfigFromJSON(jsonKey, scopes...)
	if err != nil {
		return nil, err
	}

	return c, nil
}
