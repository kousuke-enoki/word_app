package user

import (
	"context"
	"fmt"

	"word_app/backend/ent"
	"word_app/backend/ent/user"

	entsql "entgo.io/ent/dialect/sql"
)

func (r *EntUserRepo) LockByID(ctx context.Context, tx *ent.Tx, userID int) error {
	_, err := tx.User.
		Query().
		Where(user.ID(userID)).
		Modify(func(s *entsql.Selector) { s.ForUpdate() }).
		Only(ctx)
	if ent.IsNotFound(err) {
		return fmt.Errorf("user not found")
	}
	return err
}
