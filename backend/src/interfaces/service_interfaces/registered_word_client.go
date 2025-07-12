package service_interfaces

import (
	"word_app/backend/ent"
)

type RegisteredWordClient interface {
	RegisteredWord() *ent.RegisteredWordClient
}
