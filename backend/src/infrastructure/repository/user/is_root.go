package user

import (
	"context"

	"word_app/backend/ent/user"
	"word_app/backend/src/infrastructure/repoerr"
)

func (r *EntUserRepo) IsRoot(ctx context.Context, id int) (bool, error) {
	u, err := r.client.User().Query().Where(user.ID(id)).Select(user.FieldIsRoot).Only(ctx)
	if err != nil {
		return false, repoerr.FromEnt(err, "user not found", "duplicate id")
	}
	return u.IsRoot, nil
}
