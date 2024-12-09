package api

import (
	"fmt"
	db "main/db/sqlc"
	"main/token"
	"main/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server for all HTTP requests in the simple bank app
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer to setup HTTP routes and server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	router := gin.Default()
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %v", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		router:     router,
		tokenMaker: tokenMaker,
	}

	if valid, ok := binding.Validator.Engine().(*validator.Validate); ok {
		valid.RegisterValidation("currency", validCurrency)
	}

	server.setupRoutes()

	return server, nil
}

// setupRoutes sets up HTTP routes for the server
func (server *Server) setupRoutes() {

	server.router.POST("/users", server.createUser)
	server.router.POST("/users/login", server.loginUser)

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
