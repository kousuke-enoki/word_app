package config

// type LineOAuthConfig struct {
// 	ClientID     string
// 	ClientSecret string
// 	RedirectURI  string
// }

// func LoadLineConfig() LineOAuthConfig {
// 	clientID := os.Getenv("LINE_CLIENT_ID")
// 	clientSecret := os.Getenv("LINE_CLIENT_SECRET")
// 	redirectURI := os.Getenv("LINE_REDIRECT_URI")
// 	if clientID == "" || clientSecret == "" || redirectURI == "" {
// 		logrus.Fatal("LINE_CONFIG is not set")
// 	}

// 	lineCfg := LineOAuthConfig{
// 		ClientID:     clientSecret,
// 		ClientSecret: clientSecret,
// 		RedirectURI:  redirectURI,
// 	}

// 	logrus.Infof("LINE_CONFIG is set")
// 	return lineCfg
// }
