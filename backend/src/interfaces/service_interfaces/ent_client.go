package service_interfaces

import (
	"context"
	"word_app/backend/ent"
)

type EntClientInterface interface {
	Tx(ctx context.Context) (*ent.Tx, error)
	Word() *ent.WordClient
	User() *ent.UserClient
	UserConfig() *ent.UserConfigClient
	RootConfig() *ent.RootConfigClient
	RegisteredWord() *ent.RegisteredWordClient
	WordInfo() *ent.WordInfoClient
	JapaneseMean() *ent.JapaneseMeanClient
	EntClient() *ent.Client
}
