package auth

import (
	"net/http"
	"word_app/backend/src/interfaces/http/auth"
	"word_app/backend/src/utils/oauthutil"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthUsecase  auth.AuthUsecase
	jwtGenerator auth.JWTGenerator
}

func NewAuthHandler(
	authUsecase auth.AuthUsecase,
	jwtGen auth.JWTGenerator,
) *AuthHandler {
	return &AuthHandler{
		AuthUsecase:  authUsecase,
		jwtGenerator: jwtGen,
	}
}

func (h *AuthHandler) LineLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		state, _ := oauthutil.NewState(c)
		nonce, _ := oauthutil.NewNonce(c)

		url := h.AuthUsecase.StartLogin(c, state, nonce)
		c.Redirect(http.StatusFound, url)
	}
}

func (h *AuthHandler) LineCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("state")
		nonce := oauthutil.LoadNonce(c)

		res, err := h.AuthUsecase.HandleCallback(c, code, state, nonce)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func (h *AuthHandler) LineComplete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			TempToken string `json:"temp_token"`
			Password  string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		jwt, err := h.AuthUsecase.CompleteSignUp(c, req.TempToken, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": jwt})
	}
}
