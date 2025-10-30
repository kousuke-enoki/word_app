package models

type MeResponse struct {
	User User `json:"user" binding:"required"`
}

type Principal struct {
	UserID  int
	IsAdmin bool
	IsRoot  bool
	IsTest  bool
	// 将来の拡張: Scopes []string, TenantID, DeviceID など
}
