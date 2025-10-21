// src/handlers/user/edit_handler.go
package user

import (
	"net/http"
	"strconv"
	"strings"

	"word_app/backend/src/handlers/httperr"
	"word_app/backend/src/middleware/jwt"
	"word_app/backend/src/usecase/apperror"
	user_usecase "word_app/backend/src/usecase/user"
	user_validator "word_app/backend/src/validators/user"

	"github.com/gin-gonic/gin"
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

func (h *UserHandler) EditHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// 0) ルートレベルの軽い認可（最終判断はUsecase側）ついでにuserID取得
		principal, ok := jwt.GetPrincipal(c)
		if !ok || principal.IsTest {
			httperr.Write(c, apperror.Unauthorizedf("unauthorized", nil))
			return
		}

		// 1) リクエストparse
		in, err := h.parseUpdateUserRequest(c, principal.UserID)
		if err != nil {
			// parse/bind は Validation として返す
			httperr.Write(c, apperror.Validationf("invalid request", err))
			return
		}

		// 2) 追加バリデーション（フォーム系）
		if vErrs := user_validator.ValidateUpdate(*in); len(vErrs) > 0 {
			httperr.Write(c, apperror.WithFieldErrors(apperror.Validation, "invalid input", vErrs))
			return
		}

		// 3) ユースケース呼び出し（認可・検証・更新はここで）
		user, svcErr := h.userUsecase.UpdateUser(ctx, *in)
		if svcErr != nil {
			httperr.Write(c, svcErr) // ← Usecase/Repoのapperrorをそのまま
			return
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
func (h *UserHandler) parseUpdateUserRequest(c *gin.Context, userID int) (*user_usecase.UpdateUserInput, error) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	// URLパスから targetID を取得
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil || targetID <= 0 {
		return nil, apperror.Validationf("invalid id", err)
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
	in := &user_usecase.UpdateUserInput{
		EditorID:        userID,
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
