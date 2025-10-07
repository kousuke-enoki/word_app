// infra/entrepo/user_repo.go
package setting

import (
	"context"
)

func (r *EntUserConfigRepo) CreateDefault(ctx context.Context, userID int) error {
	_, err := r.client.UserConfig().
		Create().
		SetUserID(userID).
		SetIsDarkMode(false).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}
