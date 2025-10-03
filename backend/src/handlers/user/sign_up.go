package user

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/interfaces/http/user"
	"word_app/backend/src/usecase/apperror"
	user_validator "word_app/backend/src/validators/user"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignUpUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) SignUpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := h.parseRequest(c)
		if err != nil {
			httperr.Write(c, err)
			return
		}

		validationErrors := user_validator.ValidateSignUp(*req)
		if len(validationErrors) > 0 {
			httperr.Write(c, apperror.WithFieldErrors(apperror.Validation, "invalid input", validationErrors))
			return
		}

		req.Password, err = h.hashPassword(req.Password)
		if err != nil {
			httperr.Write(c, apperror.Validationf("invalid request", err))
			return
		}

		// ユーザー作成
		user, err := h.userUsecase.SignUp(context.Background(), *req)
		if err != nil {
			httperr.Write(c, err)
			return
		}

		// 作成したユーザーでサインイン（トークン発行）
		token, err := h.jwtGenerator.GenerateJWT(fmt.Sprintf("%d", user.UserID))
		if err != nil {
			httperr.Write(c, apperror.Validationf("Failed to generate token", err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Authentication successful", "token": token})
	}
}

func (h *Handler) parseRequest(c *gin.Context) (*user.SignUpInput, error) {
	if c.Request.Body == nil {
		return nil, errors.New("request body is nil")
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, errors.New("failed to read request body")
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var req SignUpUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("invalid request: " + err.Error())
	}

	// service 入力 DTO に詰め替え
	in := &user.SignUpInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	return in, nil
}

func (h *Handler) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
