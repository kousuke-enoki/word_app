// app/usecase/user/update_user.go
package user

import (
	"context"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"word_app/backend/src/app/apperr"
	"word_app/backend/src/domain"
)

type UpdateUserInput struct {
	EditorID        int
	TargetID        int
	Name            *string
	Email           *string
	PasswordNew     *string
	PasswordCurrent *string
	Role            *string // "admin" | "user" | nil=変更なし
}

type UpdateUserOutput struct {
	OK bool
}

// type UpdateUserUsecase struct {
// 	Txm      repository.TxManager
// 	UserRepo repository.UserRepository // FindByID, FindForUpdate, UpdatePartial
// }

// ざっくりメール検証（正規化はlower+trim）
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func (uc *UserUsecase) UpdateUser(ctx context.Context, in UpdateUserInput) (*UpdateUserOutput, error) {
	// Tx開始（join対応の実装前提）
	txCtx, done, err := uc.txm.Begin(ctx)
	if err != nil {
		return nil, apperr.ErrDBFailure
	}
	commit := false
	defer func() { _ = done(commit) }()

	// 1) Editor/Target取得（認可とパス検証用）
	editor, err := uc.userRepo.FindForUpdate(txCtx, in.EditorID)
	if err != nil {
		// NotFound でも基本 unauthorized
		return nil, apperr.ErrUnauthorized
	}
	target, err := uc.UserRepo.FindForUpdate(txCtx, in.TargetID)
	if err != nil {
		return nil, err // ErrUserNotFound or ErrDBFailure
	}

	// 2) 認可
	// root以外は自分のみ。testは自分でも不可
	if !editor.IsRoot && editor.ID != target.ID {
		return nil, apperr.ErrUnauthorized
	}
	if editor.ID == target.ID && editor.IsTest {
		return nil, apperr.ErrUnauthorized
	}

	// 3) 入力正規化・検証
	update := &domain.UserUpdateFields{}

	if in.Name != nil {
		n := strings.TrimSpace(*in.Name)
		update.Name = &n
	}

	if in.Email != nil {
		e := strings.ToLower(strings.TrimSpace(*in.Email))
		if !emailRegex.MatchString(e) {
			return nil, apperr.ErrInvalidParam // 400へマップする種別を用意
		}
		update.Email = &e
	}

	// パスワード更新（自分の変更 かつ 既にハッシュ有り の場合 current 必須）
	if in.PasswordNew != nil {
		if needCurrentPassword(editor, target) {
			if in.PasswordCurrent == nil {
				return nil, apperr.ErrInvalidParam // current必須
			}
			if err := verifyCurrentPassword(target.Password, *in.PasswordCurrent); err != nil {
				return nil, apperr.ErrInvalidCredential
			}
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*in.PasswordNew), bcrypt.DefaultCost)
		if err != nil {
			return nil, apperr.ErrInternal
		}
		h := string(hash)
		update.PasswordHash = &h
	}

	// 役割変更
	if in.Role != nil {
		if !editor.IsRoot {
			return nil, apperr.ErrUnauthorized
		}
		// 対象が root/test は不可
		if target.IsRoot || target.IsTest {
			return nil, apperr.ErrUnauthorized
		}
		switch strings.ToLower(*in.Role) {
		case "admin":
			b := true
			update.SetAdmin = &b
		case "user":
			b := false
			update.SetAdmin = &b
		default:
			return nil, apperr.ErrInvalidParam
		}
	}

	// 4) 更新実行
	if err := uc.UserRepo.UpdatePartial(txCtx, target.ID, update); err != nil {
		return nil, err // ErrDuplicateEmail/ErrDBFailure 等
	}

	commit = true
	if err := done(commit); err != nil {
		return nil, apperr.ErrDBFailure
	}
	return &UpdateUserOutput{OK: true}, nil
}

func needCurrentPassword(editor, target *domain.User) bool {
	if editor.ID != target.ID {
		return false
	}
	return target.Password != "" // 既にハッシュあり
}

func verifyCurrentPassword(hashed string, current string) error {
	if hashed == "" {
		return nil
	}
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(current))
}
