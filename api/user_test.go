package api

import (
	"testing"

	db "github.com/Yogksai/simplebank/db/sqlc"
	"github.com/Yogksai/simplebank/util"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:     util.RandomOwner(),
		PasswordHash: hashedPassword,
		FullName:     util.RandomOwner(),
		Email:        util.RandomEmail(),
	}
	return
}
