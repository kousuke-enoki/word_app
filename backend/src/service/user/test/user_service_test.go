package user_service_test

import (
	"context"
	"testing"

	"word_app/backend/ent/enttest"
	user_service "word_app/backend/src/service/user"

	"github.com/stretchr/testify/assert"
	_ "github.com/mattn/go-sqlite3"
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

func TestEntUserClient_FindUserByEmail(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	usrClient := user_service.NewEntUserClient(client)
	ctx := context.Background()

	// Create a user to find
	email := "test@example.com"
	name := "Test User"
	password := "password"
	_, err := usrClient.CreateUser(ctx, email, name, password)
	assert.NoError(t, err)

	// Attempt to find the user
	foundUser, err := usrClient.FindUserByEmail(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, email, foundUser.Email)
	assert.Equal(t, name, foundUser.Name)
}

func TestEntUserClient_FindUserByID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	usrClient := user_service.NewEntUserClient(client)
	ctx := context.Background()

	// Create a user to find
	email := "test@example.com"
	name := "Test User"
	password := "password"
	createdUser, err := usrClient.CreateUser(ctx, email, name, password)
	assert.NoError(t, err)

	// Attempt to find the user by ID
	foundUser, err := usrClient.FindUserByID(ctx, createdUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, email, foundUser.Email)
	assert.Equal(t, name, foundUser.Name)
}
