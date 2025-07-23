package user

import (
	"regexp"
	"unicode"

	"word_app/backend/src/models"
)

// validateSignUp checks name, email, and password fields and returns a slice of FieldError.
func ValidateSignUp(req *models.SignUpRequest) []*models.FieldError {
	var fieldErrors []*models.FieldError

	// 各フィールドの検証を個別の関数に分割
	fieldErrors = append(fieldErrors, validateUserName(req.Name)...)
	fieldErrors = append(fieldErrors, validateUserEmail(req.Email)...)
	fieldErrors = append(fieldErrors, validateUserPassword(req.Password)...)

	return fieldErrors
}

func ValidateSignIn(req *models.SignInRequest) []*models.FieldError {
	var fieldErrors []*models.FieldError

	// 各フィールドの検証を個別の関数に分割
	fieldErrors = append(fieldErrors, validateUserEmail(req.Email)...)
	fieldErrors = append(fieldErrors, validateUserPassword(req.Password)...)

	return fieldErrors
}

// Name: 長さは3〜20文字かチェック。
func validateUserName(name string) []*models.FieldError {
	var fieldErrors []*models.FieldError
	if len(name) < 3 || len(name) > 20 {
		fieldErrors = append(fieldErrors, &models.FieldError{Field: "name", Message: "name must be between 3 and 20 characters"})
	}
	return fieldErrors
}

// Email: 有効なメールアドレス形式かをチェック。
func validateUserEmail(email string) []*models.FieldError {
	var fieldErrors []*models.FieldError
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		fieldErrors = append(fieldErrors, &models.FieldError{Field: "email", Message: "invalid email format"})
	}
	return fieldErrors
}

// Password: 最低8文字、数字、アルファベットの大文字・小文字、特殊文字を含むかを確認。
func validateUserPassword(password string) []*models.FieldError {
	var fieldErrors []*models.FieldError
	if len(password) < 8 || len(password) > 30 {
		fieldErrors = append(fieldErrors, &models.FieldError{Field: "password", Message: "password must be between 8 and 30 characters"})
	}
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsNumber(ch):
			hasNumber = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}
	if !(hasUpper && hasLower && hasNumber && hasSpecial) {
		fieldErrors = append(fieldErrors, &models.FieldError{Field: "password", Message: "password must include at least one uppercase letter, one lowercase letter, one number, and one special character"})
	}
	return fieldErrors
}
