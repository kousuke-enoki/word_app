package models

type SignInRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignUpRequest struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID      int    `json:"id" binding:"required"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"isAdmin"`
	IsRoot  bool   `json:"isRoot"`
}

type MyPageResponse struct {
	User User `json:"user" binding:"required"`
}
