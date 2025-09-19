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
	IsTest  bool   `json:"isTest"`
	Email   string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	IsSettedPassword bool `json:"isSettedPassword,omitempty"`
	IsLine bool `json:"isLine,omitempty"`
}

type ExternalAuth struct {
    Provider       string `json:"provider"`
    ProviderUserId string `json:"providerUserId"`
}

type MyPageResponse struct {
	User    User `json:"user" binding:"required"`
	IsLogin bool `json:"isLogin"`
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
	Users      []User `json:"users"`
	TotalPages int    `json:"totalPages"`
}