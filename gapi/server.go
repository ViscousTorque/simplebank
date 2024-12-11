package gapi

// There are no routes in gRPC and we are not using gin with it's helpers and func

import (
	"fmt"

	db "simplebank/db/sqlc"
	"simplebank/pb"
	"simplebank/token"
	"simplebank/util"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer // this is a hack to satisfy the gRPC interface
	config                           util.Config
	store                            db.Store
	tokenMaker                       token.Maker
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
