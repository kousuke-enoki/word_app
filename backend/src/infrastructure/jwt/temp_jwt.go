package jwt

import (
	"time"
	"word_app/backend/src/usecase/auth"
	"word_app/backend/src/utils/tempjwt"
)

type TempJWTAdapter struct {
	inner *tempjwt.TempJWT
}

func New(secret string) auth.TempTokenGenerator {
	return &TempJWTAdapter{inner: tempjwt.TempJWTNew(secret)}
}

func (t *TempJWTAdapter) GenerateTemp(id *tempjwt.Identity, ttl time.Duration) (string, error) {
	return t.inner.GenerateTemp((*tempjwt.Identity)(id), ttl)
}

func (t *TempJWTAdapter) ParseTemp(tok string) (*tempjwt.Identity, error) {
	x, err := t.inner.ParseTemp(tok)
	return (*tempjwt.Identity)(x), err
}
