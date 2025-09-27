// src/service/user/update_user_service.go
package user

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"word_app/backend/ent"
	"word_app/backend/src/models"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrValidation = errors.New("validation error")
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
type FieldErrors []FieldError

func (e FieldErrors) Error() string { return ErrValidation.Error() }

// ============ 公開メソッド（薄くする） ============

func (e *EntUserClient) Update(ctx context.Context, in *models.UpdateUserInput) (*ent.User, error) {
	tx, err := e.beginTx(ctx)
	if err != nil {
		return nil, ErrDatabaseFailure
	}
	defer rollbackIfOpen(&tx)

	editor, target, err := loadEditorAndTarget(ctx, tx, in.UserID, in.TargetID)
	if err != nil {
		return nil, err
	}
	if err := authorizeUpdate(editor, target); err != nil {
		return nil, err
	}
	if err := validateUpdateInput(in); err != nil {
		return nil, err
	}

	u := buildUserUpdater(tx, target, in, editor)
	saved, err := u.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, ErrDuplicateEmail // 例: email UNIQUE
		}
		return nil, ErrDatabaseFailure
	}

	if err := tx.Commit(); err != nil {
		return nil, ErrDatabaseFailure
	}
	tx = nil
	return saved, nil
}

// ============ Tx / 共通ユーティリティ ============

func (e *EntUserClient) beginTx(ctx context.Context) (*ent.Tx, error) {
	return e.client.Tx(ctx)
}

func rollbackIfOpen(tx **ent.Tx) {
	if *tx != nil {
		_ = (*tx).Rollback()
	}
}

// ============ 取得・権限・検証 ============

func loadEditorAndTarget(ctx context.Context, tx *ent.Tx, editorID, targetID int) (*ent.User, *ent.User, error) {
	editor, err := tx.User.Get(ctx, editorID)
	if err != nil {
		return nil, nil, ErrDatabaseFailure
	}
	target, err := tx.User.Get(ctx, targetID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil, ErrUserNotFound
		}
		return nil, nil, ErrDatabaseFailure
	}
	return editor, target, nil
}

func authorizeUpdate(editor, target *ent.User) error {
	// root以外は自分のみ
	if !editor.IsRoot && editor.ID != target.ID {
		return nilOr(ErrUnauthorized)
	}
	// testユーザーは自分の編集も不可
	if editor.ID == target.ID && editor.IsTest {
		return nilOr(ErrUnauthorized)
	}
	return nil
}

// helper: 可読性のための小関数（将来ログ等を入れるならここ）
func nilOr(err error) error { return err }

func validateUpdateInput(in *models.UpdateUserInput) error {
	var verr FieldErrors
	if in == nil {
		return FieldErrors{{Field: "input", Message: "required"}}
	}
	if in.Email != nil {
		normalized := strings.ToLower(strings.TrimSpace(*in.Email))
		if !emailRegex.MatchString(normalized) {
			verr = append(verr, FieldError{Field: "email", Message: "invalid email format"})
		} else {
			in.Email = &normalized
		}
	}
	if len(verr) > 0 {
		return verr
	}
	return nil
}

// ============ 更新ビルダー（項目ごとに小関数） ============

func buildUserUpdater(tx *ent.Tx, target *ent.User, in *models.UpdateUserInput, editor *ent.User) *ent.UserUpdateOne {
	u := tx.User.UpdateOneID(target.ID)
	applyNameEmail(u, in)
	_ = applyPassword(u, editor, target, in) // エラーは Update() 側で検証済みだが、将来拡張に備えて戻り値にしておく
	_ = applyRole(u, editor, target, in)
	return u
}

func applyNameEmail(u *ent.UserUpdateOne, in *models.UpdateUserInput) {
	if in.Name != nil {
		u.SetName(strings.TrimSpace(*in.Name))
	}
	if in.Email != nil {
		u.SetEmail(*in.Email)
	}
}

func applyPassword(u *ent.UserUpdateOne, editor, target *ent.User, in *models.UpdateUserInput) error {
	if in.PasswordNew == nil {
		return nil
	}
	if needCurrentPassword(editor, target) {
		if in.PasswordCurrent == nil {
			return FieldErrors{{Field: "password.current", Message: "current password required"}}
		}
		if err := verifyCurrentPassword(target, *in.PasswordCurrent); err != nil {
			return err
		}
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(*in.PasswordNew), bcrypt.DefaultCost)
	u.SetPassword(string(hash))
	return nil
}

func needCurrentPassword(editor, target *ent.User) bool {
	// 自分自身の変更 かつ 既にパスワード設定がある場合のみ current 必須
	if editor.ID != target.ID {
		return false
	}
	return target.Password != nil && *target.Password != ""
}

func verifyCurrentPassword(target *ent.User, current string) error {
	if target.Password == nil || *target.Password == "" {
		return nil
	}
	if bcrypt.CompareHashAndPassword([]byte(*target.Password), []byte(current)) != nil {
		return FieldErrors{{Field: "password.current", Message: "current password mismatch"}}
	}
	return nil
}

func applyRole(u *ent.UserUpdateOne, editor, target *ent.User, in *models.UpdateUserInput) error {
	if in.Role == nil {
		return nil
	}
	// role変更はrootのみ
	if !editor.IsRoot {
		return ErrUnauthorized
	}
	// 対象が root/test の場合は不可
	if target.IsRoot || target.IsTest {
		return ErrUnauthorized
	}
	switch *in.Role {
	case "admin":
		u.SetIsAdmin(true)
	case "user":
		u.SetIsAdmin(false)
	default:
		return FieldErrors{{Field: "role", Message: "role must be 'admin' or 'user'"}}
	}
	return nil
}
