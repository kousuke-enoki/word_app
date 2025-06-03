package jwt

import (
	"time"
	"word_app/backend/src/interfaces/usecase/port/auth"
	"word_app/backend/src/utils/tempjwt"
)

type TempJWTAdapter struct {
	inner *tempjwt.TempJWT
}

func New(secret string) auth.TempTokenGenerator {
	return &TempJWTAdapter{inner: tempjwt.TempJWTNew(secret)}
}

func (t *TempJWTAdapter) GenerateTemp(id *auth.Identity, ttl time.Duration) (string, error) {
	return t.inner.GenerateTemp((*auth.Identity)(id), ttl)
}
func (t *TempJWTAdapter) ParseTemp(tok string) (*auth.Identity, error) {
	x, err := t.inner.ParseTemp(tok)
	return (*auth.Identity)(x), err
}
