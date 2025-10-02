// app/usecase/shared/ucerr/ucerr.go
package ucerr

import "word_app/backend/src/usecase/apperror"

func Unauthorized(msg string) error      { return apperror.Unauthorizedf(msg, nil) }
func Forbidden(msg string) error         { return apperror.Forbiddenf(msg, nil) }
func NotFound(msg string) error          { return apperror.NotFoundf(msg, nil) }
func Conflict(msg string) error          { return apperror.Conflictf(msg, nil) }
func InvalidCredential(msg string) error { return apperror.InvalidCredentialf(msg, nil) }
func Validation(msg string) error        { return apperror.Validationf(msg, nil) }
func Internal(msg string, cause error) error {
	return apperror.Internalf(msg, cause)
}

// FieldErrors（フォーム向け）
type FieldError = apperror.FieldError

func ValidationFields(msg string, fields []FieldError) error {
	return apperror.WithFieldErrors(apperror.Validation, msg, fields)
}
