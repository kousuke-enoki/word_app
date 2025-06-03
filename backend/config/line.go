package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

type LineOAuth struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func LoadLineConfig() LineOAuth {
	clientID := os.Getenv("LINE_CLIENT_ID")
	clientSecret := os.Getenv("LINE_CLIENT_SECRET")
	redirectURI := os.Getenv("LINE_REDIRECT_URI")
	if clientID == "" || clientSecret == "" || redirectURI == "" {
		logrus.Fatal("LINE_CONFIG is not set")
	}

	lineCfg := LineOAuth{
		ClientID:     clientSecret,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}

	logrus.Infof("LINE_CONFIG is set")
	return lineCfg
}
