// domain/external_auth.go
package domain

type ExternalAuth struct {
	ID             int
	UserID         int
	Provider       string // "line"
	ProviderUserID string // sub
}

func NewExternalAuth(userID int, provider, sub string) *ExternalAuth {
	return &ExternalAuth{
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: sub,
	}
}
