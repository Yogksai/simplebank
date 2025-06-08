package api

import (
	db "github.com/Yogksai/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// server serves HTTP requests for banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}
