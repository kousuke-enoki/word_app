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
	ID               int     `json:"id" binding:"required"`
	Name             string  `json:"name"`
	IsAdmin          bool    `json:"isAdmin"`
	IsRoot           bool    `json:"isRoot"`
	IsTest           bool    `json:"isTest"`
	Email            *string `json:"email,omitempty"`
	Password         string  `json:"password,omitempty"`
	IsSettedPassword bool    `json:"isSettedPassword,omitempty"`
	IsLine           bool    `json:"isLine,omitempty"`
	CreatedAt        string  `json:"createdAt,omitempty"`
	UpdatedAt        string  `json:"updatedAt,omitempty"`
}

type UserDetail struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Email            *string `json:"email,omitempty"`
	IsAdmin          bool    `json:"isAdmin"`
	IsRoot           bool    `json:"isRoot"`
	IsTest           bool    `json:"isTest"`
	IsLine           bool    `json:"isLine"`
	IsSettedPassword bool    `json:"isSettedPassword"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}

type ExternalAuth struct {
	Provider       string `json:"provider"`
	ProviderUserId string `json:"providerUserId"`
}

type MyPageResponse struct {
	User    User `json:"user" binding:"required"`
	IsLogin bool `json:"isLogin"`
}
