package user

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"word_app/backend/src/models"
	user_service "word_app/backend/src/service/user"
	"word_app/backend/src/validators/user"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

		logrus.Info("j")
		validationErrors := user.ValidateSignUp(req)
		if len(validationErrors) > 0 {
			logrus.Info("validationErrors")
			logrus.Info(validationErrors)
			c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
			return
		}
		logrus.Info("adsf")

		hashedPassword, err := h.hashPassword(req.Password)
		if err != nil {
			logrus.Info("err")
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		logrus.Info("qwer")

		// ユーザー作成
		user, err := h.userClient.CreateUser(context.Background(), req.Email, req.Name, hashedPassword)
		if err != nil {
			logrus.Error(err)

			// エラーの種類ごとにレスポンスを変更
			switch err {
			case user_service.ErrDuplicateEmail:
				c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			case user_service.ErrDatabaseFailure:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "An unknown error occurred"})
			}
			return
		}
		logrus.Info("ert")

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
