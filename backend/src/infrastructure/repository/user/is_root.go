package user

import (
	"context"

	"word_app/backend/ent/user"
)

func (r *EntUserRepo) IsRoot(ctx context.Context, id int) (bool, error) {
	u, err := r.client.User().Query().Where(user.ID(id)).Select(user.FieldIsRoot).Only(ctx)
	if err != nil {
		return false, err
	}
	return u.IsRoot, nil
}
