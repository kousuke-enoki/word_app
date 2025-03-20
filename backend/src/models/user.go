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
	Name  string `json:"name" binding:"required"`
	Admin bool   `json:"admin" binding:"required"`
	Root  bool   `json:"root" binding:"required"`
}

type MyPageResponse struct {
	User User `json:"user" binding:"required"`
}
