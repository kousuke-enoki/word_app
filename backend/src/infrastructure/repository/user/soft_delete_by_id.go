package user

import (
	"context"
	"time"

	"word_app/backend/src/infrastructure/repoerr"
)

func (r *EntUserRepo) SoftDeleteByID(ctx context.Context, id int, t time.Time) error {
	_, err := r.client.User().
		UpdateOneID(id).
		SetDeletedAt(t).
		Save(ctx)
	if err != nil {
		return repoerr.FromEnt(err, "internal", "internal server error")
	}
	return nil
}
