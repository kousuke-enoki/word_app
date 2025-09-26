package user

import (
	"context"
	"time"
)

func (r *EntUserRepo) SoftDeleteByID(ctx context.Context, id int, t time.Time) error {
	_, err := r.client.User().
		UpdateOneID(id).
		SetDeletedAt(t).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
}
