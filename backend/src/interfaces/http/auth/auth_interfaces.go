package auth

import (
	"context"

	"github.com/gin-gonic/gin"
)

type CallbackResult struct {
	Token         string `json:"token,omitempty"`
	NeedPassword  bool   `json:"need_password,omitempty"`
	TempToken     string `json:"temp_token,omitempty"`
	SuggestedMail string `json:"suggested_mail,omitempty"`
}

type Handler interface {
	LineLogin() gin.HandlerFunc
	LineCallback() gin.HandlerFunc
	LineComplete() gin.HandlerFunc
}

type JWTGenerator interface {
	GenerateJWT(userID string) (string, error)
}

type AuthUsecase interface {
	StartLogin(ctx context.Context, state, nonce string) string
	HandleCallback(ctx context.Context, code, state, nonce string) (*CallbackResult, error)
	CompleteSignUp(ctx context.Context, tempToken, pass string) (string, error)
}
