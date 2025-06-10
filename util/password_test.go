package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	err = CheckPasswordHash(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPasswordHash(wrongPassword, hashedPassword)
	require.EqualError(t, err, fmt.Errorf("password does not match hash:  %w", bcrypt.ErrMismatchedHashAndPassword).Error())
}
