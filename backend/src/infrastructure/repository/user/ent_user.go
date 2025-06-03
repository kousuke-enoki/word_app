// infrastructure/repository/user_ent.go
package user

import (
	"context"

	"word_app/backend/ent/externalauth"
	"word_app/backend/ent/user"
	"word_app/backend/src/domain"
	"word_app/backend/src/interfaces/service_interfaces"

	"github.com/sirupsen/logrus"
)

type EntUserRepo struct {
	client service_interfaces.EntClientInterface
}

func NewEntUserRepo(c service_interfaces.EntClientInterface) *EntUserRepo {
	return &EntUserRepo{client: c}
}

func (r *EntUserRepo) FindByProvider(ctx context.Context, provider, subject string) (*domain.User, error) {
	u, err := r.client.User().
		Query().
		Where(
			user.HasExternalAuthsWith(
				externalauth.Provider(provider),
				externalauth.ProviderUserID(subject),
			)).
		Only(ctx)
	if err != nil {
		return nil, err // ent.ErrNotFound なら呼び出し側で nil 判定
	}
	return &domain.User{ID: u.ID, Email: u.Email, Name: u.Name}, nil
}

func (r *EntUserRepo) Create(ctx context.Context, u *domain.User, ext *domain.ExternalAuth) error {
	// トランザクション開始
	tx, err := r.client.Tx(ctx)
	if err != nil {
		logrus.Error(err)
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				logrus.Error(err)
			}
		}
	}()

	eu, err := tx.User.
		Create().
		SetEmail(u.Email).
		SetName(u.Name).
		SetPassword(u.Password).
		Save(ctx)
	if err != nil {
		return err
	}
	if _, err = tx.ExternalAuth.
		Create().
		SetUserID(eu.ID).
		SetProvider(ext.Provider).
		SetProviderUserID(ext.ProviderUserID).
		Save(ctx); err != nil {
		return err
	}
	return nil
}
