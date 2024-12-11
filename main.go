package main

import (
	"context"
	"log"
	"net"
	"os"

	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/gapi"
	"simplebank/pb"
	"simplebank/util"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the db: ", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	runGrpcServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create a new server", err)
	}

	err = server.Run(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot run a new server", err)
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	_ = os.Setenv("GRPC_GO_LOG_LEVEL", "debug")
	_ = os.Setenv("GRPC_GO_LOG_SEVERITY_LEVEL", "debug")
	_ = os.Setenv("GRPC_TRACE", "all")
	_ = os.Setenv("GRPC_VERBOSITY_LEVEL", "debug")

	logger := log.Default() // This is the standard Go logger, you can replace it with logrus or zap
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(logger.Writer(), logger.Writer(), logger.Writer()))

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create a new server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	// Reflection allows you to interactively query the service using tools like Evans.
	// This should be removed in production to avoid exposing your service details.
	if config.EnableReflection {
		log.Println("Enabling gRPC reflection...")
		reflection.Register(grpcServer)
	} else {
		log.Println("Reflection is disabled for security reasons.")
	}

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}
	log.Printf("start gRPC listener server at %s", config.GrpcServerAddress)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("gRPC server failed to serve")
	}
}
