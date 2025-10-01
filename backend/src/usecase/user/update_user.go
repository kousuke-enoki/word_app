// app/usecase/user/update_user.go
package user

import (
	"context"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"word_app/backend/src/domain"
	"word_app/backend/src/domain/repository"
	"word_app/backend/src/interfaces/http/user"
	"word_app/backend/src/models"
	"word_app/backend/src/usecase/apperror"
)

// ざっくりメール検証（正規化はlower+trim）
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func (uc *UserUsecase) UpdateUser(ctx context.Context, in user.UpdateUserInput) (*models.UserDetail, error) {
	// 1) Tx
	txCtx, done, err := uc.txm.Begin(ctx)
	if err != nil {
		return nil, err
	}
	commit := false
	defer func() { _ = done(commit) }()

	// 2) 取得
	editor, target, err := uc.loadEditorAndTarget(txCtx, in.EditorID, in.TargetID)
	if err != nil {
		return nil, err
	}

	// 3) 認可
	if err := uc.authorizeUpdate(editor, target); err != nil {
		return nil, err
	}

	// 4) 入力→更新フィールド
	update, err := uc.buildUpdateFields(editor, target, in)
	if err != nil {
		return nil, err
	}

	// 5) 更新
	updated, err := uc.userRepo.UpdatePartial(txCtx, target.ID, update)
	if err != nil {
		return nil, err
	}

	// 6) Commit & DTO
	commit = true
	if err := done(commit); err != nil {
		return nil, err
	}
	return uc.toUserDetail(updated), nil
}

// ---------- helpers (小さな関数へ分割) ----------

func (uc *UserUsecase) loadEditorAndTarget(ctx context.Context, editorID, targetID int) (*domain.User, *domain.User, error) {
	editor, err := uc.userRepo.FindForUpdate(ctx, editorID)
	if err != nil {
		return nil, nil, err
	}
	target, err := uc.userRepo.FindForUpdate(ctx, targetID)
	if err != nil {
		return nil, nil, err
	}
	return editor, target, nil
}

func (uc *UserUsecase) authorizeUpdate(editor, target *domain.User) error {
	// root以外は自分のみ
	if !editor.IsRoot && editor.ID != target.ID {
		return uc.errUnauthorized("Unauthorized", nil)
	}
	// testは自分でも不可
	if editor.ID == target.ID && editor.IsTest {
		return uc.errUnauthorized("Unauthorized", nil)
	}
	return nil
}

// 入力正規化・検証 → repository.UserUpdateFields を構築
func (uc *UserUsecase) buildUpdateFields(editor, target *domain.User, in user.UpdateUserInput) (*repository.UserUpdateFields, error) {
	out := &repository.UserUpdateFields{}

	// name
	if in.Name != nil {
		n := strings.TrimSpace(*in.Name)
		out.Name = &n
	}

	// email
	if in.Email != nil {
		e, err := uc.normalizeEmail(*in.Email)
		if err != nil {
			return nil, err
		}
		out.Email = &e
	}

	// password
	if in.PasswordNew != nil {
		hash, err := uc.hashNewPasswordIfAllowed(editor, target, in.PasswordCurrent, *in.PasswordNew)
		if err != nil {
			return nil, err
		}
		out.PasswordHash = &hash
	}

	// role
	if in.Role != nil {
		admin, err := uc.resolveRoleChange(editor, target, *in.Role)
		if err != nil {
			return nil, err
		}
		out.SetAdmin = &admin
	}

	return out, nil
}

func (uc *UserUsecase) normalizeEmail(raw string) (string, error) {
	e := strings.ToLower(strings.TrimSpace(raw))
	if !emailRegex.MatchString(e) {
		return "", uc.errInvalidParam("VALIDATION", nil)
	}
	return e, nil
}

func (uc *UserUsecase) hashNewPasswordIfAllowed(editor, target *domain.User, currentOpt *string, newPlain string) (string, error) {
	if uc.needCurrentPassword(editor, target) {
		if currentOpt == nil {
			return "", uc.errInvalidParam("VALIDATION", nil) // current必須
		}
		if err := uc.verifyCurrentPassword(target.Password, *currentOpt); err != nil {
			return "", uc.errInvalidCredential("ERR_INVALID_CREDENTIAL", nil)
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPlain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (uc *UserUsecase) resolveRoleChange(editor, target *domain.User, role string) (bool, error) {
	// role変更はrootのみ / 対象がroot/testは不可
	if !editor.IsRoot || target.IsRoot || target.IsTest {
		return false, uc.errUnauthorized("Unauthorized", nil)
	}
	switch strings.ToLower(role) {
	case "admin":
		return true, nil
	case "user":
		return false, nil
	default:
		return false, uc.errInvalidParam("VALIDATION", nil)
	}
}

func (uc *UserUsecase) needCurrentPassword(editor, target *domain.User) bool {
	return editor.ID == target.ID && target.Password != ""
}

func (uc *UserUsecase) verifyCurrentPassword(hashed, current string) error {
	if hashed == "" {
		return nil
	}
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(current))
}

func (uc *UserUsecase) toUserDetail(u *domain.User) *models.UserDetail {
	var emailPtr *string
	if u.Email != nil {
		e := *u.Email
		emailPtr = &e
	}
	return &models.UserDetail{
		ID:      u.ID,
		Name:    u.Name,
		Email:   emailPtr,
		IsAdmin: u.IsAdmin,
		IsRoot:  u.IsRoot,
		IsTest:  u.IsTest,
		// 必要に応じて追加
	}
}

// ---- エラー種別を一箇所に（適宜あなたのapperrに差し替えOK） ----

func (uc *UserUsecase) errUnauthorized(msg string, err *error) error {
	return apperror.New("UNAUTHORIZED", msg, *err)
} // 例: apperr.ErrUnauthorized
func (uc *UserUsecase) errInvalidParam(msg string, err *error) error {
	return apperror.New("VALIDATION", msg, *err)
} // 例: apperr.ErrInvalidParam
func (uc *UserUsecase) errInvalidCredential(msg string, err *error) error {
	return apperror.New("ERR_INVALID_CREDENTIAL", msg, *err)
} // 例: apperr.ErrInvalidCredential
