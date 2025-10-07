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
	var emailPtr *string
	if email != nil { // Ent も Nillable にした前提
		email := *email   // string 取り出し
		emailPtr = &email // ポインタ化（そのまま u.Email でも良い）
	}
	var passPtr *string
	if rawPass != nil { // Ent も Nillable にした前提
		rawPass := *rawPass // string 取り出し
		passPtr = &rawPass  // ポインタ化（そのまま u.Email でも良い）
	}
	hash, err := hashPassword(passPtr)
	if err != nil {
		return nil, err
	}
	return &User{
		Email:    emailPtr,
		Name:     name,
		Password: hash,
	}, nil
}

func hashPassword(pass *string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(*pass), bcrypt.DefaultCost)
	return string(b), err
}
