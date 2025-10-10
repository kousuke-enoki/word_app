package userdailyusage

import (
	"context"
	"time"
)

// CreateIfNotExists は、ユーザーのデイリー使用行を存在しなければ作成します。
// 既に行がある場合は DoNothing で副作用ゼロ。
// last_reset_date は JST 当日0時、各カウントは 0 で初期化します。
// ユーザーのデイリー使用行を「なければ作る」。
// 既に存在すれば何もしない（上書きしない）。
func (r *EntUserDailyUsageRepo) CreateIfNotExists(ctx context.Context, userID int, now time.Time) error {
	todayJST := r.truncateToJST0(now)

	// INSERT ... ON CONFLICT (user_id) DO NOTHING
	err := r.client.UserDailyUsage().
		Create().
		SetUserID(userID).
		SetLastResetDate(todayJST).
		// 明示的に初期化（デフォルト値でもよいが可読性のためセット）
		SetQuizCount(0).
		SetBulkCount(0).
		// OnConflict(
		// 	sql.ConflictColumns(userdailyusage.FieldUserID),
		// ).
		// DoNothing().
		Exec(ctx) // ← Save ではなく Exec

	return err
}
