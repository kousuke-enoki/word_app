// domain/user.go
package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Email    *string
	Name     string
	Password string
	IsRoot   bool
	IsAdmin  bool
	IsTest   bool

	// これらは Domain の関心。UI ではない
	HasPassword bool
	HasLine     bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name string, email, rawPass *string) (*User, error) {
	u := &User{
		Name:    name,
		IsAdmin: false, IsRoot: false, IsTest: false,
	}
	if email != nil { // Ent も Nillable にした前提
		email := *email  // string 取り出し
		u.Email = &email // ポインタ化（そのまま u.Email でも良い）
	}
	if rawPass == nil {
		// パスワード無しユーザー（外部認証 or テストユーザー）
		u.Password = ""
		return u, nil
	}
	var passPtr *string
	pass := *rawPass // string 取り出し
	passPtr = &pass  // ポインタ化（そのまま u.Email でも良い）
	hash, err := hashPassword(passPtr)
	if err != nil {
		return nil, err
	}
	u.Password = hash
	u.HasPassword = rawPass != nil
	u.HasLine = false // 後で紐付けるまで false
	return u, nil
}

func hashPassword(pass *string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(*pass), bcrypt.DefaultCost)
	return string(b), err
}
