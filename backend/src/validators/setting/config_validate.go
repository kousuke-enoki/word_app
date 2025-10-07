package setting

import (
	"word_app/backend/src/models"
	settingUc "word_app/backend/src/usecase/setting"
)

type Validator interface {
	ValidateRootConfig(SignUpRequest *settingUc.InputUpdateRootConfig) []*models.FieldError
}

func ValidateRootConfig(req *settingUc.InputUpdateRootConfig) []*models.FieldError {
	var fieldErrors []*models.FieldError

	// 各フィールドの検証を個別の関数に分割
	fieldErrors = append(fieldErrors, validateEditingPermissions(req.EditingPermission)...)

	return fieldErrors
}

func validateEditingPermissions(editingPermissions string) []*models.FieldError {
	var fieldErrors []*models.FieldError
	validRoles := map[string]bool{
		"user":  true,
		"admin": true,
		"root":  true,
	}
	if !validRoles[editingPermissions] {
		fieldErrors = append(fieldErrors, &models.FieldError{
			Field:   "editing_permissions",
			Message: "editing_permissions must be one of: user, admin, root",
		})
	}
	return fieldErrors
}
