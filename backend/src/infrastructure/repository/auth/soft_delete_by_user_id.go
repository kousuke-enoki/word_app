package auth

import (
	"context"
	"time"

	"word_app/backend/ent/externalauth"
	"word_app/backend/ent/user"
)

func (r *EntExtAuthRepo) SoftDeleteByUserID(ctx context.Context, userID int, t time.Time) error {
	_, err := r.client.ExternalAuth().
		Update().
		Where(
			externalauth.HasUserWith(user.ID(userID)),
			externalauth.DeletedAtIsNil(),
		).
		SetDeletedAt(t).
		Save(ctx)
	return err
}
