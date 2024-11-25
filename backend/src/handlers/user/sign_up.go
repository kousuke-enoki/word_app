package user

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"unicode"

	"word_app/backend/src/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// FieldError defines an error structure for specific fields
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (h *UserHandler) SignUpHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := h.parseRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErrors := h.validateSignUp(req)
		if len(validationErrors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
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

		token, err := h.jwtGenerator.GenerateJWT(fmt.Sprintf("%d", user.ID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Authentication successful", "token": token})
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

// validateSignUp checks name, email, and password fields and returns a slice of FieldError.
func (h *UserHandler) validateSignUp(req *models.SignUpRequest) []FieldError {
	// var errors []FieldError
	var fieldErrors []FieldError

	// 各フィールドの検証を個別の関数に分割
	fieldErrors = append(fieldErrors, h.validateName(req.Name)...)
	fieldErrors = append(fieldErrors, h.validateEmail(req.Email)...)
	fieldErrors = append(fieldErrors, h.validatePassword(req.Password)...)

	return fieldErrors
}

// Name: 長さは3〜20文字かチェック。
func (h *UserHandler) validateName(name string) []FieldError {
	var fieldErrors []FieldError
	if len(name) < 3 || len(name) > 20 {
		fieldErrors = append(fieldErrors, FieldError{Field: "name", Message: "name must be between 3 and 20 characters"})
	}
	return fieldErrors
}

// Email: 有効なメールアドレス形式かをチェック。
func (h *UserHandler) validateEmail(email string) []FieldError {
	var fieldErrors []FieldError
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		fieldErrors = append(fieldErrors, FieldError{Field: "email", Message: "invalid email format"})
	}
	return fieldErrors
}

// Password: 最低8文字、数字、アルファベットの大文字・小文字、特殊文字を含むかを確認。
func (h *UserHandler) validatePassword(password string) []FieldError {
	var fieldErrors []FieldError
	if len(password) < 8 || len(password) > 30 {
		fieldErrors = append(fieldErrors, FieldError{Field: "password", Message: "password must be between 8 and 30 characters"})
	}
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsNumber(ch):
			hasNumber = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}
	if !(hasUpper && hasLower && hasNumber && hasSpecial) {
		fieldErrors = append(fieldErrors, FieldError{Field: "password", Message: "password must include at least one uppercase letter, one lowercase letter, one number, and one special character"})
	}
	return fieldErrors
}
