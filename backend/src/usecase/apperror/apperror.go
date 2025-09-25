// usecase/apperror/apperror.go
package apperror

type Kind string

const (
	Unauthorized Kind = "UNAUTHORIZED"
	Forbidden    Kind = "FORBIDDEN"
	NotFound     Kind = "NOT_FOUND"
	Conflict     Kind = "CONFLICT"
	Validation   Kind = "VALIDATION"
	Internal     Kind = "INTERNAL"
)

type Error struct {
	Kind    Kind   // 意味コード
	Message string // ユーザーに見せてもよい短文
	Err     error  // 内部原因（wrap用）
}

func (e *Error) Error() string { return e.Message }

func New(kind Kind, msg string, err error) *Error {
	return &Error{Kind: kind, Message: msg, Err: err}
}
