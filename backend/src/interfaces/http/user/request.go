package user

type UpdateUserInput struct {
	EditorID        int
	TargetID        int
	Name            *string
	Email           *string
	PasswordNew     *string
	PasswordCurrent *string
	Role            *string // "admin" | "user" | nil=変更なし
}
type SignUpInput struct {
	Email    string
	Name     string
	Password string
}
type ListUsersInput struct {
	ViewerID int
	Search   string
	SortBy   string // "name" | "email" | "role"
	Order    string // "asc" | "desc"
	Page     int
	Limit    int
}
type DeleteUserInput struct {
	EditorID int // 操作者
	TargetID int // 削除対象
}
