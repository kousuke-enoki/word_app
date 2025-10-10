// usecase/apperror/apperror.go
package apperror

import (
	"errors"

	"word_app/backend/src/models"
)

type Kind string

const (
	Unauthorized      Kind = "UNAUTHORIZED"
	Forbidden         Kind = "FORBIDDEN"
	NotFound          Kind = "NOT_FOUND"
	Conflict          Kind = "CONFLICT"
	Validation        Kind = "VALIDATION"
	Internal          Kind = "INTERNAL"
	InvalidCredential Kind = "INVALID_CREDENTIAL"
	TooManyRequests   Kind = "TOO_MANY_REQUESTS"
)

type Error struct {
	Kind    Kind           // 分類（HandlerでHTTPへマップ）
	Message string         // ユーザー向け短文（i18nキーでもOK）
	Err     error          // 内部原因（ログ用）
	Meta    map[string]any // 追加情報（field errors 等）
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return string(e.Kind)
}

func (e *Error) Unwrap() error { return e.Err }

// コンストラクタ
func New(kind Kind, msg string, err error) *Error {
	return &Error{Kind: kind, Message: msg, Err: err}
}

func WithFieldErrors(kind Kind, msg string, fields []models.FieldError) *Error {
	return &Error{Kind: kind, Message: msg, Meta: map[string]any{"fields": fields}}
}

// ほかの処理内で使用できるショートカット
func Unauthorizedf(msg string, err error) *Error      { return New(Unauthorized, msg, err) }
func Forbiddenf(msg string, err error) *Error         { return New(Forbidden, msg, err) }
func NotFoundf(msg string, err error) *Error          { return New(NotFound, msg, err) }
func Conflictf(msg string, err error) *Error          { return New(Conflict, msg, err) }
func Validationf(msg string, err error) *Error        { return New(Validation, msg, err) }
func InvalidCredentialf(msg string, err error) *Error { return New(InvalidCredential, msg, err) }
func Internalf(msg string, err error) *Error          { return New(Internal, msg, err) }
func TooManyRequestsf(msg string, err error) *Error   { return New(TooManyRequests, msg, err) }

// ユーティリティ
func IsKind(err error, k Kind) bool {
	var ae *Error
	return errors.As(err, &ae) && ae.Kind == k
}
