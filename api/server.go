package api

import (
	db "github.com/Yogksai/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// server serves HTTP requests for banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// Constructor function to create a new server
func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()
	router.POST("/account", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
