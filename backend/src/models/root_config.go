package models

type RootConfig struct {
	EditingPermission string `json:"editing_permission" binding:"required,oneof=user admin root"`
	IsTestUserMode    bool   `json:"is_test_user_mode"`
	IsEmailAuthCheck  bool   `json:"is_email_authentication_check"`
	IsLineAuth        bool   `json:"is_line_authentication"`
}
