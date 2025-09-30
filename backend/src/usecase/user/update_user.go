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
)

// ざっくりメール検証（正規化はlower+trim）
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func (uc *UserUsecase) UpdateUser(ctx context.Context, in user.UpdateUserInput) (*models.UserDetail, error) {
	// Tx開始（join対応の実装前提）
	txCtx, done, err := uc.txm.Begin(ctx)
	if err != nil {
		return nil, err
	}
	commit := false
	defer func() { _ = done(commit) }()

	// 1) Editor/Target取得（認可とパス検証用）
	editor, err := uc.userRepo.FindForUpdate(txCtx, in.EditorID)
	if err != nil {
		// NotFound でも基本 unauthorized
		return nil, err
	}
	target, err := uc.userRepo.FindForUpdate(txCtx, in.TargetID)
	if err != nil {
		return nil, err // ErrUserNotFound or ErrDBFailure
	}

	// 2) 認可
	// root以外は自分のみ。testは自分でも不可
	if !editor.IsRoot && editor.ID != target.ID {
		return nil, err
	}
	if editor.ID == target.ID && editor.IsTest {
		return nil, err
	}

	// 3) 入力正規化・検証
	update := &repository.UserUpdateFields{}

	if in.Name != nil {
		n := strings.TrimSpace(*in.Name)
		update.Name = &n
	}

	if in.Email != nil {
		e := strings.ToLower(strings.TrimSpace(*in.Email))
		if !emailRegex.MatchString(e) {
			return nil, err // 400へマップする種別を用意
		}
		update.Email = &e
	}

	// パスワード更新（自分の変更 かつ 既にハッシュ有り の場合 current 必須）
	if in.PasswordNew != nil {
		if needCurrentPassword(editor, target) {
			if in.PasswordCurrent == nil {
				return nil, err // current必須
			}
			if err := verifyCurrentPassword(target.Password, *in.PasswordCurrent); err != nil {
				return nil, err
			}
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*in.PasswordNew), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		h := string(hash)
		update.PasswordHash = &h
	}

	// 役割変更
	if in.Role != nil {
		if !editor.IsRoot {
			return nil, err
		}
		// 対象が root/test は不可
		if target.IsRoot || target.IsTest {
			return nil, err
		}
		switch strings.ToLower(*in.Role) {
		case "admin":
			b := true
			update.SetAdmin = &b
		case "user":
			b := false
			update.SetAdmin = &b
		default:
			return nil, err
		}
	}

	// 4) 更新実行
	updatedUser, err := uc.userRepo.UpdatePartial(txCtx, target.ID, update)
	if err != nil {
		return nil, err // ErrDuplicateEmail/ErrDBFailure 等
	}

	commit = true
	if err := done(commit); err != nil {
		return nil, err
	}
	return &models.UserDetail{
		ID:      updatedUser.ID,
		Name:    updatedUser.Name,
		Email:   updatedUser.Email,
		IsAdmin: updatedUser.IsAdmin,
		IsRoot:  updatedUser.IsRoot,
		IsTest:  updatedUser.IsTest,
		// Add other fields as needed from updatedUser
	}, nil
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
