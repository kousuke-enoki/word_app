package user

import (
	"word_app/backend/src/interfaces/http/user"
	"word_app/backend/src/models"
	userfields "word_app/backend/src/validators/user/userFields"
)

// validateSignUp checks name, email, and password fields and returns a slice of FieldError.
func ValidateSignUp(req *user.SignUpInput) []*models.FieldError {
	var fieldErrors []*models.FieldError

	// 各フィールドの検証を個別の関数に分割
	fieldErrors = append(fieldErrors, userfields.ValidateUserName(req.Name)...)
	fieldErrors = append(fieldErrors, userfields.ValidateUserEmail(req.Email)...)
	fieldErrors = append(fieldErrors, userfields.ValidateUserPassword(req.Password)...)

	return fieldErrors
}

func ValidateSignIn(req *models.SignInRequest) []*models.FieldError {
	var fieldErrors []*models.FieldError

	// 各フィールドの検証を個別の関数に分割
	fieldErrors = append(fieldErrors, userfields.ValidateUserEmail(req.Email)...)
	fieldErrors = append(fieldErrors, userfields.ValidateUserPassword(req.Password)...)

	return fieldErrors
}
