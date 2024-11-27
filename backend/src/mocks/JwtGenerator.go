package mocks

import (
	mock "github.com/stretchr/testify/mock"
)

// JWT トークン生成のモック関数
type JwtGenerator struct {
	mock.Mock
}

func (m *JwtGenerator) GenerateJWT(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}
