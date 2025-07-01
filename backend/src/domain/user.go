// domain/user.go
package domain

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       int
	Email    string
	Name     string
	Password string
	IsRoot   bool
	IsAdmin  bool
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
