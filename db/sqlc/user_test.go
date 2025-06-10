package db

import (
	"context"
	"testing"

	"github.com/Yogksai/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err, "Error getting user")
	require.NotEmpty(t, user2, "User should not be empty")

	require.Equal(t, user1.Username, user2.Username, "Usernames do not match")
	require.Equal(t, user1.PasswordHash, user2.PasswordHash, "Password hashes do not match")
	require.Equal(t, user1.FullName, user2.FullName, "Full names do not match")
	require.Equal(t, user1.Email, user2.Email, "Emails do not match")
	require.Equal(t, user1.PasswordChangedAt.Time.IsZero(), user2.PasswordChangedAt.Time.IsZero(), "Password changed at times do not match")
	require.NotZero(t, user2.CreatedAt)
}

func CreateRandomUser(t *testing.T) User {
	HashPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err, "Error hashing password")
	require.NotEmpty(t, HashPassword, "Hashed password should not be empty")
	arg := CreateUserParams{
		Username:     util.RandomOwner(),
		PasswordHash: HashPassword,
		FullName:     util.RandomOwner(),
		Email:        util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.Time.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}
