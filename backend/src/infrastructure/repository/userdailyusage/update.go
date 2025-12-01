// src/infrastructure/repository/ent_user_daily_usage_updater.go
package userdailyusage

import (
	"context"
	"database/sql"
	"time"

	"word_app/backend/src/domain"
	"word_app/backend/src/infrastructure/repoerr"
	"word_app/backend/src/usecase/apperror"
)

type UsageKind int

const (
	UsageQuiz UsageKind = iota + 1
	UsageBulk
)

// --- public ---

func (r *EntUserDailyUsageRepo) IncQuizOr429(ctx context.Context, userID int, now time.Time, dailyCap int) (*domain.DailyUsageUpdateResult, error) {
	return r.incWithKind(ctx, userID, now, UsageQuiz, dailyCap)
}

func (r *EntUserDailyUsageRepo) IncBulkOr429(ctx context.Context, userID int, now time.Time, dailyCap int) (*domain.DailyUsageUpdateResult, error) {
	return r.incWithKind(ctx, userID, now, UsageBulk, dailyCap)
}

// --- internal ---

func (r *EntUserDailyUsageRepo) incWithKind(ctx context.Context, userID int, now time.Time, kind UsageKind, dailyCap int) (*domain.DailyUsageUpdateResult, error) {
	if dailyCap <= 0 {
		// ガード（設定漏れに強く）
		if kind == UsageQuiz {
			dailyCap = 20
		} else {
			dailyCap = 5
		}
	}
	today := r.truncateToJST0(now)

	var q string
	switch kind {
	case UsageQuiz:
		// 昨日なら quiz=1, bulk=0 にリセット / 今日なら quiz<cap のときだけ +1
		q = quizAddSql()
	case UsageBulk:
		// 昨日なら bulk=1, quiz=0 にリセット / 今日なら bulk<cap のときだけ +1
		q = bulkAddSql()
	default:
		return nil, apperror.Internalf("unknown usage kind", nil)
	}

	var res domain.DailyUsageUpdateResult
	// 注意：ent.Client.DB() は *sql.DB を返す前提（プロジェクトのセットアップに依存）
	err := r.sql.
		QueryRowContext(ctx, q, userID, today, dailyCap).
		Scan(&res.QuizCount, &res.BulkCount, &res.LastResetDate)
	if err != nil {
		// SQLのWHERE句で更新されなかった場合（上限到達）は sql.ErrNoRows が返る
		// これを429エラーに変換
		if err == sql.ErrNoRows {
			return nil, apperror.TooManyRequestsf("daily quota exceeded", nil)
		}
		// データベースエラーをラップ
		return nil, repoerr.FromEnt(err, "failed to update daily usage", "database error")
	}

	// 上限判定：更新後カウントが cap のままで、今回インクリメントされていないケースを検出するには
	// 直前値が必要になりますが、UPSERT+CASEでは“増えなかった”＝“すでにcap”と等価。
	// 呼び出し側の意図に合わせ、「capを超えたらエラー（429）」を返す設計にします。
	switch kind {
	case UsageQuiz:
		if res.QuizCount > dailyCap {
			// 理論上発生しないがガード
			return &res, apperror.TooManyRequestsf("daily quota exceeded", nil)
		}
		// 「capに到達」自体は成功（20回目はOK）。21回目の呼び出しはres.QuizCount==capのままなのでエラーにしたい場合、
		// 直前値と比較が必要。簡易には、ハンドラ側で「この呼び出しの前に既にcapか」を覚えておくか、
		// ここで “cap到達後の再呼び出し” を弾くための追加SELECTを入れる。
		// 実務では「capに到達していたら以降エラー」で良いケースが多いので、簡易に以下を採用：
		// if res.QuizCount == cap { // 到達した瞬間もOKにしたいならここは通過
		// 	// ここではOKとして返す（20回目成功）。21回目でまたこの値が返るので呼び出し側でエラー化しても良い。
		// }
	case UsageBulk:
		if res.BulkCount > dailyCap {
			return &res, apperror.TooManyRequestsf("daily quota exceeded", nil)
		}
		// if res.BulkCount == cap {
		// 	// 同上のコメント。5回目はOK、6回目は値が増えないので呼び出し側で429にするのがシンプル。
		// }
	}

	return &res, nil
}

// func (r *EntUserDailyUsageRepo) truncateToJST0(t time.Time) time.Time {
// 	tt := t.In(r.jst)
// 	y, m, d := tt.Date()
// 	return time.Date(y, m, d, 0, 0, 0, 0, r.jst)
// }

func quizAddSql() string {
	return `
		INSERT INTO user_daily_usages (user_id, last_reset_date, quiz_count, bulk_count, updated_at)
		VALUES ($1, $2, 1, 0, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id)
		DO UPDATE SET
			quiz_count = CASE
				WHEN user_daily_usages.last_reset_date < $2 THEN 1
				ELSE user_daily_usages.quiz_count + 1
			END,
			bulk_count = CASE
				WHEN user_daily_usages.last_reset_date < $2 THEN 0
				ELSE user_daily_usages.bulk_count
			END,
			last_reset_date = CASE
				WHEN user_daily_usages.last_reset_date < $2 THEN $2
				ELSE user_daily_usages.last_reset_date
			END,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_daily_usages.last_reset_date < $2
		   OR user_daily_usages.quiz_count < $3
		RETURNING quiz_count, bulk_count, last_reset_date;
	`
}

func bulkAddSql() string {
	return `
		INSERT INTO user_daily_usages (user_id, last_reset_date, quiz_count, bulk_count, updated_at)
		VALUES ($1, $2, 0, 1, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id)
		DO UPDATE	SET
			bulk_count = CASE
				WHEN user_daily_usages.last_reset_date < $2 THEN 1
				ELSE user_daily_usages.bulk_count + 1
			END,
			quiz_count = CASE
				WHEN user_daily_usages.last_reset_date < $2 THEN 0
				ELSE user_daily_usages.quiz_count
			END,
			last_reset_date = CASE
				WHEN user_daily_usages.last_reset_date < $2 THEN $2
				ELSE user_daily_usages.last_reset_date
			END,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_daily_usages.last_reset_date < $2
			OR user_daily_usages.bulk_count < $3
		RETURNING quiz_count, bulk_count, last_reset_date;
	`
}
