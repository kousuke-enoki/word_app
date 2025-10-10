// domain/user.go
package domain

import (
	"time"
)

type UserDailyUsage struct {
	ID            int
	LastResetDate time.Time
	QuizCount     int
	BulkCount     int
	UpdatedAt     time.Time
}

type DailyUsageUpdateResult struct {
	QuizCount     int
	BulkCount     int
	LastResetDate time.Time
}

// func NewUserDailyUsage(userID int, QuizCount, BulkCount *int) (*UserDailyUsage, error) {
// 	todayJST := r.truncateToJST0(now)
// 	hash, err := hashPassword(passPtr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &User{
// 		Email:    emailPtr,
// 		Name:     name,
// 		Password: hash,
// 	}, nil
// }

// // JSTのその日0:00に丸める
// func truncateToJST0(t time.Time) time.Time {
// 	tt := t.In(r.jst)
// 	y, m, d := tt.Date()
// 	return time.Date(y, m, d, 0, 0, 0, 0, r.jst)
// }
