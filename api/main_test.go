package api

import (
	db "github.com/dangquyit/go-simplebank/db/sqlc"
	"github.com/dangquyit/go-simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		AccessTokenDuration: time.Minute,
		TokenSymmetricKey:   util.RandomString(32),
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
