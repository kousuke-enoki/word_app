package auth

import (
	"os"

	"golang.org/x/oauth2"
)

var lineEndpoint = oauth2.Endpoint{
	AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
	TokenURL: "https://api.line.me/oauth2/v2.1/token", // :contentReference[oaicite:0]{index=0}
}

func (h *AuthHandler) oauthCfg() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("LINE_CLIENT_ID"),
		ClientSecret: os.Getenv("LINE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("LINE_REDIRECT_URI"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     lineEndpoint,
	}
}
