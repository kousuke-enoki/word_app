package setting

import (
	"context"
	"time"

	"word_app/backend/ent/userconfig"
)

func (r *EntUserConfigRepo) SoftDeleteByUserID(ctx context.Context, userID int, t time.Time) error {
	_, err := r.client.UserConfig().
		Update().
		Where(userconfig.UserID(userID), userconfig.DeletedAtIsNil()).
		SetDeletedAt(t).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}
