package api

import (
	"os"
	"testing"
	"time"

	db "github.com/Yogksai/simplebank/db/sqlc"
	"github.com/Yogksai/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenKey:            util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(store, config)
	require.NoError(t, err)
	require.NotEmpty(t, server)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
