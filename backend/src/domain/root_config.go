package domain

import "time"

type RootConfig struct {
	ID                         int
	EditingPermission          string    `json:"editing_permission"`
	IsTestUserMode             bool      `json:"is_test_user_mode"`
	IsEmailAuthenticationCheck bool      `json:"is_email_authentication_check"`
	IsLineAuthentication       bool      `json:"is_line_authentication"`
	UpdatedAt                  time.Time `json:"updated_at"`
}
