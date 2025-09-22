// src/handlers/user/edit_handler.go
package user

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"word_app/backend/src/models"
	user_service "word_app/backend/src/service/user"
	"word_app/backend/src/utils/contextutil"
	user_validator "word_app/backend/src/validators/user"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UpdatePasswordPayload struct {
	New     *string `json:"new"`
	Current *string `json:"current"`
}

type UpdateUserRequest struct {
	Name     *string                `json:"name"`
	Email    *string                `json:"email"`
	Password *UpdatePasswordPayload `json:"password"`
	Role     *string                `json:"role"` // "admin" | "user"
}

func (h *Handler) EditHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		userRoles, err := contextutil.GetUserRoles(c)
		if err != nil || userRoles == nil || userRoles.IsTest {
			if err == nil {
				err = errors.New("unauthorized: admin access required")
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// リクエストを解析
		in, err := h.parseUpdateUserRequest(c)
		if err != nil {
			logrus.Errorf("Failed to parse request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// バリデーション
		if vErrs := user_validator.ValidateUpdate(in); len(vErrs) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"errors": vErrs})
			return
		}

		// サービス呼び出し
		user, svcErr := h.userClient.Update(ctx, in)
		if svcErr != nil {
			// エラー種別に応じて HTTP へマッピング
			switch {
			case errors.Is(svcErr, user_service.ErrUnauthorized):
				c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			case errors.Is(svcErr, user_service.ErrUserNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			case errors.Is(svcErr, user_service.ErrDuplicateEmail):
				c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
				return
			case errors.Is(svcErr, user_service.ErrValidation):
				// サービスから詳細バリデーションを返す場合
				if fe, ok := svcErr.(user_service.FieldErrors); ok {
					c.JSON(http.StatusBadRequest, gin.H{"errors": fe})
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "validation error"})
				}
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
				return
			}
		}

		// 更新後のサマリーを最小限返す
		c.JSON(http.StatusOK, gin.H{
			"id":      user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"isAdmin": user.IsAdmin,
			"isRoot":  user.IsRoot,
			"isTest":  user.IsTest,
		})
	}
}

// リクエスト構造体を解析
func (h *Handler) parseUpdateUserRequest(c *gin.Context) (*models.UpdateUserInput, error) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	// ユーザーIDをコンテキストから取得
	userID, exists := c.Get("userID")
	if !exists {
		return nil, errors.New("unauthorized: userID not found in context")
	}

	// userIDの型チェック
	userIDInt, ok := userID.(int)
	if !ok {
		return nil, errors.New("invalid userID type")
	}

	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil || targetID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return nil, errors.New("invalid id")
	}

	// 追加の軽い正規化（email は service 側でも最終正規化）
	if req.Email != nil {
		e := strings.TrimSpace(*req.Email)
		req.Email = &e
	}
	if req.Name != nil {
		n := strings.TrimSpace(*req.Name)
		req.Name = &n
	}
	if req.Password != nil && req.Password.New != nil {
		p := strings.TrimSpace(*req.Password.New)
		req.Password.New = &p
	}
	if req.Password != nil && req.Password.Current != nil {
		p := strings.TrimSpace(*req.Password.Current)
		req.Password.Current = &p
	}

	// service 入力 DTO に詰め替え
	in := &models.UpdateUserInput{
		UserID:          userIDInt,
		TargetID:        targetID,
		Name:            req.Name,
		Email:           req.Email,
		PasswordNew:     nil,
		PasswordCurrent: nil,
		Role:            req.Role,
	}
	if req.Password != nil {
		in.PasswordNew = req.Password.New
		in.PasswordCurrent = req.Password.Current
	}

	return in, nil
}
