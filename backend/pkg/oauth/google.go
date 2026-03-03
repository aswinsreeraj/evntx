package oauth

import (
	"context"
	"os"

	"google.golang.org/api/idtoken"
)

type GoogleUser struct {
	Email         string
	EmailVerified bool
	Name          string
}

func VerifyGoogleIDToken(idToken string) (*GoogleUser, error) {
	payload, err := idtoken.Validate(context.Background(), idToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return nil, err
	}

	return &GoogleUser{
		Email:         payload.Claims["email"].(string),
		EmailVerified: payload.Claims["email_verified"].(bool),
		Name:          payload.Claims["name"].(string),
	}, nil
}
