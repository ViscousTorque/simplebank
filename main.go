package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"

	"simplebank/api"
	db "simplebank/db/sqlc"
	_ "simplebank/doc/statik"
	"simplebank/gapi"
	"simplebank/mail"
	"simplebank/pb"
	"simplebank/util"
	"simplebank/worker"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		var notFoundErr *viper.ConfigFileNotFoundError
		if errors.As(err, &notFoundErr) {
			log.Warn().Msgf("Config file not found, but continuing with defaults: %v", notFoundErr)
		} else {
			log.Fatal().Err(err).Msg("Cannot load config")
		}
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to the db")
	}
	defer conn.Close()

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	redisOpts := asynq.RedisClientOpt{ // we are not include prod TLS server conn for secure conn
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpts)

	//TODO: use wait group here for better shutdown handling!
	go runTaskProcessor(config, redisOpts, store) // go routine, because redis will block and keep polling for new tasks
	go runGrpcServer(config, store, taskDistributor)
	runGatewayServer(config, store, taskDistributor)
	// runGinServer(config, store)
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)

	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}

}

func runDBMigration(migrationURL string, dbSource string) {
	log.Info().Msg("Starting db migration ...")
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	// its not clear to me why there would be an error return if there is no change!
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

// TODO: setup this on config or something :-)
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create a new server")
	}

	err = server.Run(config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot run a new server")
	}
}

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create a new server")
	}

	gprcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(gprcLogger)

	pb.RegisterSimpleBankServer(grpcServer, server)
	if config.EnableReflection {
		log.Info().Msg("Enabling gRPC reflection...")
		reflection.Register(grpcServer)
	} else {
		log.Fatal().Err(err).Msg("Reflection is disabled for security reasons.")
	}

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}
	log.Info().Msgf("start gRPC listener server at %s", config.GrpcServerAddress)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("gRPC server failed to serve")
	}
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create a new server")
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
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create http listener")
	}
	log.Info().Msgf("start Http Gateway server at %s", listener.Addr().String())

	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("gRPC Http Gateway failed to serve")
	}
}
