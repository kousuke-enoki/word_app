package user

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"word_app/backend/src/handlers"
	user_interface "word_app/backend/src/interfaces/http/user"
	"word_app/backend/src/validators/user"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// FieldError defines an error structure for specific fields
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (h *Handler) SignUpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := h.parseRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErrors := user.ValidateSignUp(req)
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}

		req.Password, err = h.hashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// ユーザー作成
		user, err := h.userUsecase.SignUp(context.Background(), *req)
		if err != nil {
			handlers.WriteError(c, err)
			return
		}
		// if err != nil {
		// 	// エラーの種類ごとにレスポンスを変更
		// 	switch err {
		// 	case user_service.ErrDuplicateEmail:
		// 		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		// 	case user_service.ErrDatabaseFailure:
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		// 	default:
		// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unknown error occurred"})
		// 	}
		// 	return
		// }

		token, err := h.jwtGenerator.GenerateJWT(fmt.Sprintf("%d", user.UserID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Authentication successful", "token": token})
	}
}

func (h *Handler) parseRequest(c *gin.Context) (*user_interface.SignUpInput, error) {
	if c.Request.Body == nil {
		return nil, errors.New("request body is nil")
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, errors.New("failed to read request body")
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var req user_interface.SignUpInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("invalid request: " + err.Error())
	}
	return &req, nil
}

func (h *Handler) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
