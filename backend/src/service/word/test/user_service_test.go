package user_service_test

import (
	"context"
	"testing"

	"word_app/backend/ent/enttest"
	user_service "word_app/backend/src/service/user"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestEntUserClient_CreateUser(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	usrClient := user_service.NewEntUserClient(client)
	ctx := context.Background()

	email := "test@example.com"
	name := "Test User"
	password := "password"

	createdUser, err := usrClient.CreateUser(ctx, email, name, password)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, email, createdUser.Email)
	assert.Equal(t, name, createdUser.Name)
}
