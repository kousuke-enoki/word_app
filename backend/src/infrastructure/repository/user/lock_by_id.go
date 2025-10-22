package user

import (
	"context"
	"fmt"

	"word_app/backend/ent"
	"word_app/backend/ent/user"
	"word_app/backend/src/infrastructure/repository/tx"

	entsql "entgo.io/ent/dialect/sql"
)

func (r *EntUserRepo) LockByID(ctx context.Context, userID int) error {
	// txCtx に埋め込まれた *ent.Tx を取得
	txx, ok := tx.TxFromContext(ctx)
	if !ok || txx == nil {
		return fmt.Errorf("LockByID must be called within a transaction")
	}

	_, err := txx.User.
		Query().
		Where(user.ID(userID)).
		Modify(func(s *entsql.Selector) { s.ForUpdate() }). // ent v0.14 の FOR UPDATE
		Only(ctx)

	if ent.IsNotFound(err) {
		return fmt.Errorf("user not found")
	}
	return err
}
