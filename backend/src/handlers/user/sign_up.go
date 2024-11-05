package user

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"word_app/backend/src/models"
	"word_app/backend/src/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *UserHandler) SignUpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := h.parseRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := h.hashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user, err := h.userClient.CreateUser(context.Background(), req.Email, req.Name, hashedPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Sign up failed", "details": err.Error()})
			return
		}

		// サインアップ後にJWTトークンを生成
		token, err := utils.GenerateJWT(fmt.Sprintf("%d", user.ID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func (h *UserHandler) parseRequest(c *gin.Context) (*models.SignUpRequest, error) {
	if c.Request.Body == nil {
		return nil, errors.New("request body is nil")
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, errors.New("failed to read request body")
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var req models.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("invalid request: " + err.Error())
	}
	return &req, nil
}

func (h *UserHandler) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
