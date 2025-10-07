package user

import "word_app/backend/src/models"

type SignUpOutput struct {
	UserID int
}

type UserListResponse struct {
	Users      []models.User `json:"users"`
	TotalPages int           `json:"totalPages"`
}
type FindByEmailOutput struct {
	UserID         int
	HashedPassword string
	IsAdmin        bool
	IsRoot         bool
	IsTest         bool
	// 必要ならEmailやNameも返せる
}
