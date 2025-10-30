package auth

import (
	"net/http"

	"word_app/backend/src/utils/oauthutil"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) LineLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		state, _ := oauthutil.NewState(c)
		nonce, _ := oauthutil.NewNonce(c)

		ctx := c.Request.Context()

		url := h.AuthUsecase.StartLogin(ctx, state, nonce)
		c.Redirect(http.StatusFound, url)
	}
}

func (h *AuthHandler) LineCallback() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		code := c.Query("code")
		// state := c.Query("state")
		// nonce := oauthutil.LoadNonce(c)

		res, err := h.AuthUsecase.HandleCallback(ctx, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func (h *AuthHandler) LineComplete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var req struct {
			TempToken string  `json:"temp_token"`
			Password  *string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		jwt, err := h.AuthUsecase.CompleteSignUp(ctx, req.TempToken, req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": jwt})
	}
}
