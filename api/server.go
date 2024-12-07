package api

import (
	db "main/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server for all HTTP requests in the simple bank app
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer to setup HTTP routes and server
func NewServer(store db.Store) *Server {
	router := gin.Default()
	server := &Server{
		store:  store,
		router: router,
	}

	if valid, ok := binding.Validator.Engine().(*validator.Validate); ok {
		valid.RegisterValidation("currency", validCurrency)
	}

	server.setupRoutes()

	return server
}

// setupRoutes sets up HTTP routes for the server
func (server *Server) setupRoutes() {

	server.router.POST("/users", server.createUser)

	server.router.POST("/accounts", server.createAccount)
	server.router.GET("/accounts/:id", server.getAccount)
	server.router.GET("/accounts", server.listAccounts)

	server.router.POST("/transfers", server.createTransfers)
}

// Run starts the HTTP server
func (server *Server) Run(address string) error {
	//TODO: add graceful shutdown code
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
