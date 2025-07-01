package domain

type RootConfig struct {
	ID                         int
	EditingPermission          string
	IsTestUserMode             bool
	IsEmailAuthenticationCheck bool
	IsLineAuthentication       bool
}
