package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/glowfi/voxpopuli/backend/internal/middleware"
	postrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/post"
	postsvc "github.com/glowfi/voxpopuli/backend/pkg/service/post"
	transport "github.com/glowfi/voxpopuli/backend/pkg/transport"
	"github.com/joho/godotenv"
	"github.com/oklog/run"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	// Initialize logger
	logger := zerolog.New(os.Stdout)

	// Create a context that can be canceled
	ctx, shutdownFunc := context.WithCancel(context.Background())
	defer shutdownFunc()

	// Setup database
	databaseDSN := "postgres://%s:%s@%s/%s?sslmode=disable"
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	// Connect to the database
	dsn := fmt.Sprintf(databaseDSN, dbUsername, dbPassword, dbHost, dbName)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	if err := db.Ping(); err != nil {
		logger.Fatal().Err(err).Msg("failed to ping database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Err(err).Msg("db connection failed to close")
		}
	}()

	// Initialize repo and services
	postRepo := postrepo.NewRepo(db)
	postSvc := postsvc.NewService(postRepo)

	services := transport.Services{
		Post: postSvc,
	}

	// Create a new transportServer
	transportServer, err := transport.NewServer(services)
	if err != nil {
		logger.Fatal().Err(err).Msg("server creation failed")
	}

	// Create an HTTP handler for the server
	httpHandler, err := transportServer.HTTPHandler(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("http handler failed to start")
	}
	httpHandler = http.StripPrefix("/api", httpHandler)

	// Create a root router
	rootRouter := http.NewServeMux()

	// create middleware stack
	corsOptions := middleware.DefaultCORSOptions()
	middlewareStack := middleware.CreateStack(
		middleware.Logging,
		middleware.CORS(corsOptions),
	)
	rootRouter.Handle("/api/", middlewareStack(httpHandler))

	// Create an HTTP server
	portStr := os.Getenv("SERVER_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to cast server port from string to integer")
	}
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           rootRouter,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		BaseContext:       func(_ net.Listener) context.Context { return ctx },
	}

	var rg run.Group

	// start http server
	rg.Add(func() error {
		logger.Info().Msgf("starting http server on port %d", port)
		return httpServer.ListenAndServe()
	}, func(err error) {
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("http server exited")
		}
	})

	// graceful shutdown
	quitC := make(chan os.Signal, 1)
	rg.Add(func() error {
		signal.Notify(quitC, os.Interrupt, syscall.SIGTERM)
		sig := <-quitC
		return fmt.Errorf("received signal %s", sig)
	}, func(error) {
		shutdownFunc()
	})

	if err := rg.Run(); err != nil {
		logger.Err(err).Msg("exiting")
	}
}
