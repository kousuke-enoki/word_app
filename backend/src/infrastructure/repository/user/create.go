package user

import (
	"context"

	"word_app/backend/src/domain"

	"github.com/sirupsen/logrus"
)

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

	var emailPtr *string
	if u.Email != nil { // Ent も Nillable にした前提
		Email := *u.Email // string 取り出し
		emailPtr = &Email // ポインタ化（そのまま u.Email でも良い）
	}

	eu, err := tx.User.
		Create().
		SetNillableEmail(emailPtr).
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
