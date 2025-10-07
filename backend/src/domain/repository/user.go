package repository

import "word_app/backend/src/domain"

type UserListFilter struct {
	Search string
	SortBy string // "name" | "email" | "role"
	Order  string // "asc" | "desc"
	Offset int
	Limit  int
}

type UserListResult struct {
	Users      []*domain.User // HasPassword/HasLine を含む
	TotalCount int
}

type UserUpdateFields struct {
	Name         *string // nil=変更なし
	Email        *string // 正規化済み（lower/trim）
	PasswordHash *string // 新パスワードのハッシュ（nil=変更なし）
	// 役割（is_admin）
	SetAdmin *bool // nil=変更なし、true/falseで変更
}
