package jwt

import (
	"time"

	"word_app/backend/src/utils/tempjwt"
)

type TempJWTAdapter struct {
	inner *tempjwt.TempJWT
}

type TempTokenGenerator interface {
	GenerateTemp(id *tempjwt.Identity, ttl time.Duration) (string, error)
	ParseTemp(tok string) (*tempjwt.Identity, error)
}

func New(secret string) TempTokenGenerator {
	return &TempJWTAdapter{inner: tempjwt.New(secret)}
}

func (t *TempJWTAdapter) GenerateTemp(id *tempjwt.Identity, ttl time.Duration) (string, error) {
	return t.inner.GenerateTemp((*tempjwt.Identity)(id), ttl)
}

func (t *TempJWTAdapter) ParseTemp(tok string) (*tempjwt.Identity, error) {
	x, err := t.inner.ParseTemp(tok)
	return (*tempjwt.Identity)(x), err
}
