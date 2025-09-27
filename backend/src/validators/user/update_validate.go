// src/validators/user/update.go
package user

import (
	"word_app/backend/src/models"
	userfields "word_app/backend/src/validators/user/userFields"
)

// Update の入口バリデーション。
// - in.Name / in.Email / in.PasswordNew / in.Role のうち、指定されたもののみ検証します。
func ValidateUpdate(in *models.UpdateUserInput) []*models.FieldError {
	var errs []*models.FieldError

	// name
	if in.Name != nil {
		if e := userfields.ValidateUserName(*in.Name); len(e) > 0 {
			errs = append(errs, e...)
		}
	}

	// email
	if in.Email != nil {
		if e := userfields.ValidateUserEmail(*in.Email); len(e) > 0 {
			errs = append(errs, e...)
		}
	}

	// password.new（current の要否は DB 状態依存 → service 側で判定）
	if in.PasswordNew != nil {
		if e := userfields.ValidateUserPassword(*in.PasswordNew); len(e) > 0 {
			// SignUp と同じ「password」フィールド名で返すか、
			// 画面での区別用に "password.new" にしたい場合はここを変えてください。
			// ここでは SignUp に合わせ "password" に寄せます。
			for _, fe := range e {
				fe.Field = "password"
			}
			errs = append(errs, e...)
		}
	}

	// role（root/test の扱いは service 側で判定）
	if in.Role != nil {
		if *in.Role != "admin" && *in.Role != "user" {
			errs = append(errs, &models.FieldError{
				Field:   "role",
				Message: "role must be 'admin' or 'user'",
			})
		}
	}

	return errs
}
