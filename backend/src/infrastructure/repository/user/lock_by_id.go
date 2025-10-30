// repository/user/lock_by_id.go
package user

import (
	"context"
	"fmt"

	txRepo "word_app/backend/src/infrastructure/repository/tx"
)

func (r *EntUserRepo) LockByID(ctx context.Context, userID int) error {
	tx, ok := txRepo.TxFromContext(ctx)
	if !ok || tx == nil {
		return fmt.Errorf("no transaction in context")
	}
	// アプリ名の名前空間を混ぜておくと他用途と鍵空間が衝突しません
	// 例: 0x7770_0001 << 32 | userID（int8にパック）
	key := int64(0x77700001)<<32 | int64(userID)
	_, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock($1)`, key)
	return err
}

// 子テーブルへ INSERT するとき、PostgreSQL は親（users）の該当行に
// FOR KEY SHARE を取りに行きます（FK 整合性確認のため）。
// こちらの LockByID が親行に FOR UPDATE をかける実装だと、
// FOR KEY SHARE と FOR UPDATE は衝突するため、同一ユーザーで同時登録や、
// 別 Tx でロックを取っている状況だと待機に入ります。
// もし LockByID が 「同じ TxCtx を使っていない」（= 別コネクション／別トランザクションでロックを取っている）と、
// 自己衝突で待ち続けます。これも起きている可能性が高いです。
// func (r *EntUserRepo) LockByID(ctx context.Context, userID int) error {
// 	// txCtx に埋め込まれた *ent.Tx を取得
// 	txx, ok := tx.TxFromContext(ctx)
// 	if !ok || txx == nil {
// 		return fmt.Errorf("LockByID must be called within a transaction")
// 	}

// 	_, err := txx.User.
// 		Query().
// 		Where(user.ID(userID)).
// 		Modify(func(s *entsql.Selector) { s.ForUpdate() }). // ent v0.14 の FOR UPDATE
// 		Only(ctx)

// 	if ent.IsNotFound(err) {
// 		return fmt.Errorf("user not found")
// 	}
// 	return err
// }
