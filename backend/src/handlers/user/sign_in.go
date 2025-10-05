package user

import (
	"context"
	"fmt"
	"net/http"

	"word_app/backend/src/handlers"
	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/apperror"
	"word_app/backend/src/validators/user"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) SignInHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SignInRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			if fs := handlers.FieldsFromBindError(err); len(fs) > 0 {
				httperr.Write(c, apperror.WithFieldErrors(apperror.Validation, "invalid input", fs))
				return
			}
			// バリデータ以外のbindエラーも 400 に寄せたいならこちら
			httperr.Write(c, apperror.Validationf("invalid input", err))
			return
		}
		validationErrors := user.ValidateSignIn(req)
		if len(validationErrors) > 0 {
			httperr.Write(c, apperror.WithFieldErrors(apperror.Validation, "invalid input", validationErrors))
			return
		}

		// ユーザー検索
		signInUser, err := h.userUsecase.FindByEmail(context.Background(), req.Email)
		if err != nil {
			httperr.Write(c, err)
			return
		}

		// パスワード未設定（外部認証のみ）のユーザーはパスワードサインイン不可
		if err := comparePasswordPtr(signInUser.HashedPassword, req.Password); err != nil {
			// エラーメッセージは統一して情報リークを防ぐ
			httperr.Write(c, apperror.Validationf("invalid request", nil))
			return
		}
		// JWT発行
		token, err := h.jwtGenerator.GenerateJWT(fmt.Sprintf("%d", signInUser.UserID))
		if err != nil {
			httperr.Write(c, apperror.Validationf("invalid request", nil))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Authentication successful", "token": token,
		})
	}
}

// hashed が nil/空ならエラー。平文とbcryptで比較。
func comparePasswordPtr(hashed string, plain string) error {
	if hashed == "" {
		return apperror.Validationf("password not set", nil)
	}
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
