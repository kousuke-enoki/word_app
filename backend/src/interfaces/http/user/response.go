package user

import "word_app/backend/src/models"

type SignUpOutput struct {
	UserID int
}
type UserListRequest struct {
	UserID int    `json:"userId"`
	Search string `json:"search"`
	SortBy string `json:"sortBy"`
	Order  string `json:"order"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
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
