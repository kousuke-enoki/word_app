// src/service/user/update.go
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

// 既存のものを再利用 or 追加
var ErrValidation = errors.New("validation error")

// 任意：詳細エラーを返したいときの型
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
type FieldErrors []FieldError

func (e FieldErrors) Error() string { return ErrValidation.Error() }

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (e *EntUserClient) Update(ctx context.Context, in *models.UpdateUserInput) (*ent.User, error) {
	tx, err := e.client.Tx(ctx)
	if err != nil {
		return nil, ErrDatabaseFailure
	}
	defer func() {
		// rollback if still open
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	editor, err := tx.User.Get(ctx, in.UserID)
	if err != nil {
		return nil, ErrDatabaseFailure
	}
	target, err := tx.User.Get(ctx, in.TargetID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrUserNotFound
		}
		return nil, ErrDatabaseFailure
	}

	// 権限: root 以外は自分のみ
	if !editor.IsRoot && editor.ID != target.ID {
		return nil, ErrUnauthorized
	}
	// test ユーザーは自分の編集も不可
	if editor.ID == target.ID && editor.IsTest {
		return nil, ErrUnauthorized
	}

	// バリデーション（ここでは正規化だけ確認する。詳細は validators/user/update.go で）
	// ※ in のポインタは nil 可能なので、nil チェックを忘れずに
	// ※ ここで詳細バリデーションを行い、FieldErrors を返すことも可能
	//    例: 必須チェック、文字数チェックなど
	//    ただし、複雑になる場合は validators/user/update.go に分離した方が良い
	var verr FieldErrors

	if in.Email != nil {
		normalized := strings.ToLower(strings.TrimSpace(*in.Email))
		if !emailRegex.MatchString(normalized) {
			verr = append(verr, FieldError{Field: "email", Message: "invalid email format"})
		} else {
			in.Email = &normalized
		}
	}

	if len(verr) > 0 {
		return nil, verr
	}

	u := tx.User.UpdateOneID(target.ID)

	// 反映
	if in.Name != nil {
		u.SetName(strings.TrimSpace(*in.Name))
	}
	if in.Email != nil {
		u.SetEmail(*in.Email)
	}

	// password 更新
	if in.PasswordNew != nil {
		needCurrent := false
		if editor.ID == target.ID {
			// 自分のPW変更時、現状PWが設定済なら current 必須
			if target.Password != nil && *target.Password != "" {
				needCurrent = true
			}
		}
		if needCurrent {
			if in.PasswordCurrent == nil {
				return nil, FieldErrors{{Field: "password.current", Message: "current password required"}}
			}
			if bcrypt.CompareHashAndPassword([]byte(*target.Password), []byte(*in.PasswordCurrent)) != nil {
				return nil, FieldErrors{{Field: "password.current", Message: "current password mismatch"}}
			}
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(*in.PasswordNew), bcrypt.DefaultCost)
		hashed := string(hash)
		u.SetPassword(hashed)
	}

	// role 変更
	if in.Role != nil {
		// root だけが role を変えられる
		if !editor.IsRoot {
			return nil, ErrUnauthorized
		}
		// 対象が root/test の場合は role 変更不可
		if target.IsRoot || target.IsTest {
			return nil, ErrUnauthorized
		}
		switch *in.Role {
		case "admin":
			u.SetIsAdmin(true)
		case "user":
			u.SetIsAdmin(false)
		default:
			return nil, FieldErrors{{Field: "role", Message: "role must be 'admin' or 'user'"}}
		}
		// isRoot/isTest は触らない
	}

	saved, err := u.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			// email UNIQUE など
			return nil, ErrDuplicateEmail
		}
		return nil, ErrDatabaseFailure
	}

	if err = tx.Commit(); err != nil {
		return nil, ErrDatabaseFailure
	}
	// tx = nil で defer の Rollback を抑止
	tx = nil

	return saved, nil
}
