package domain

type UserConfig struct {
	ID         int
	UserID     int
	IsDarkMode bool `json:"is_dark_mode"` // true=ダーク・false=ライト
}
