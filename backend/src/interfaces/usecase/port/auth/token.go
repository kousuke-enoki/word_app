package auth

import (
	"time"
)

// type Identity struct {
// 	Provider string `json:"provider"`
// 	Sub      string `json:"sub"`
// 	Email    string `json:"email"`
// 	Name     string `json:"name"`
// 	jwt.RegisteredClaims
// }

type TempTokenGenerator interface {
	GenerateTemp(id *Identity, ttl time.Duration) (string, error)
	ParseTemp(tok string) (*Identity, error)
}
