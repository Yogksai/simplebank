package api

import (
	"fmt"

	db "github.com/Yogksai/simplebank/db/sqlc"
	token "github.com/Yogksai/simplebank/token"
	"github.com/Yogksai/simplebank/util"
	"github.com/gin-gonic/gin"
)

// server serves HTTP requests for banking service
type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker //interface
	config     util.Config
}

// Constructor function to create a new server
func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.InitRoutes()
	return server, nil
}

func (server *Server) InitRoutes() {
	router := gin.Default()
	router.POST("/users/login", server.loginUser)
	router.POST("/users", server.createUser)

	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRouter.POST("/accounts", server.createAccount)
	authRouter.GET("/accounts/:id", server.getAccount)
	authRouter.GET("/accounts", server.listAccounts)

	authRouter.POST("/transfers", server.createTransfer)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
