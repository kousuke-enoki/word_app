// infra/entrepo/user_repo.go
package user

import (
	"context"

	"word_app/backend/ent"
	"word_app/backend/ent/user"
	"word_app/backend/src/domain"
	"word_app/backend/src/domain/repository"
	usermapper "word_app/backend/src/infrastructure/mapper/user"
)

type UserRepository struct {
	Client *ent.Client
}

// Tx/非Tx切替の共通ヘルパ（既出）
// func getDB(ctx context.Context, client *ent.Client) interface {
// 	User() *ent.UserClient
// } {
// 	if tx, ok := txFromContext(ctx); ok && tx != nil { // 既存のTxManager実装に合わせる
// 		return tx
// 	}
// 	return client
// }

func (r *EntUserRepo) FindForUpdate(ctx context.Context, id int) (*domain.User, error) {
	u, err := r.client.User().
		Query().
		Where(user.ID(id), user.DeletedAtIsNil()).
		Select(
			user.FieldID,
			user.FieldName,
			user.FieldEmail,
			user.FieldPassword,
			user.FieldIsAdmin,
			user.FieldIsRoot,
			user.FieldIsTest,
			user.FieldCreatedAt,
			user.FieldUpdatedAt,
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, err
		}
		return nil, err
	}
	var emailPtr *string
	if u.Email != nil {
		e := *u.Email
		emailPtr = &e
	}
	var pass string
	if u.Password != nil {
		pass = *u.Password
	}
	return &domain.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     emailPtr,
		Password:  pass, // ハッシュ
		IsAdmin:   u.IsAdmin,
		IsRoot:    u.IsRoot,
		IsTest:    u.IsTest,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (r *EntUserRepo) UpdatePartial(ctx context.Context, targetID int, f *repository.UserUpdateFields) (*domain.User, error) {
	u := r.client.User().UpdateOneID(targetID)

	if f.Name != nil {
		u.SetName(*f.Name)
	}
	if f.Email != nil {
		u.SetEmail(*f.Email) // Nillableなら SetNillableEmail(f.Email) に変更
	}
	if f.PasswordHash != nil {
		u.SetPassword(*f.PasswordHash)
	}
	if f.SetAdmin != nil {
		u.SetIsAdmin(*f.SetAdmin)
	}
	user, err := u.Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, err
		}
		return nil, err
	}
	return usermapper.MapEntUser(user, nil), nil
}
