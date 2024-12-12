package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/gapi"
	"simplebank/pb"
	"simplebank/util"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	//TODO: use wait group here for better shutdown handling!
	go runGrpcServer(config, store)
	runGatewayServer(config, store)
}

// TODO: setup this on config or something :-)
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
		log.Fatal("cannot create listener", err)
	}
	log.Printf("start gRPC listener server at %s", config.GrpcServerAddress)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("gRPC server failed to serve", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	logger := log.Default() // This is the standard Go logger, you can replace it with logrus or zap
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(logger.Writer(), logger.Writer(), logger.Writer()))

	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create a new server", err)
	}

	// enable snake case as per proto files - see gateway doco
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatalln("cannot register handler server", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot create http listener", err)
	}
	log.Printf("start Http Gateway server at %s", listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("gRPC Http Gateway failed to serve", err)
	}
}
