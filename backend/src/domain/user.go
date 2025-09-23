// domain/user.go
package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Email    string
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

func NewUser(email, name, rawPass string) (*User, error) {
	hash, err := hashPassword(rawPass)
	if err != nil {
		return nil, err
	}
	return &User{
		Email:    email,
		Name:     name,
		Password: hash,
	}, nil
}

func hashPassword(pass string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(b), err
}
