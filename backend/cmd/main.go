package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/glowfi/voxpopuli/backend/internal/middleware"
	postrepo "github.com/glowfi/voxpopuli/backend/pkg/repo/post"
	postsvc "github.com/glowfi/voxpopuli/backend/pkg/service/post"
	transport "github.com/glowfi/voxpopuli/backend/pkg/transport"
	posttransport "github.com/glowfi/voxpopuli/backend/pkg/transport/post"
	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

const (
	port        = 8080
	databaseDSN = "postgres://%s:%s@%s/%s?sslmode=disable"
)

func main() {
	// Load configuration
	dbUsername := "postgres"
	dbPassword := "postgres"
	dbHost := "127.0.0.1:5432"
	dbName := "voxpopuli"

	// Initialize logger
	logger := zerolog.New(os.Stdout)

	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to the database
	dsn := fmt.Sprintf(databaseDSN, dbUsername, dbPassword, dbHost, dbName)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	if err := db.Ping(); err != nil {
		logger.Fatal().Err(err).Msg("")
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Err(err).Msg("db connection failed to close")
		}
	}()

	// Initialize services and transports
	postRepo := postrepo.NewRepo(db)
	postSvc := postsvc.NewService(postRepo)
	postsTransport := posttransport.NewTransport(postSvc)

	services := transport.Services{
		Post: postSvc,
	}
	transports := transport.Transports{
		Post: postsTransport,
	}

	// Create a new server
	server, err := transport.NewServer(services, transports)
	if err != nil {
		logger.Fatal().Err(err).Msg("server creation failed")
	}

	// Create an HTTP handler for the server
	apiHandler, err := server.HTTPHandler(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("api handler failed to start")
	}
	apiHandler = http.StripPrefix("/api", apiHandler)

	// Create a root router
	rootRouter := http.NewServeMux()

	// create middleware stack
	corsOptions := middleware.DefaultCORSOptions()
	middlewareStack := middleware.CreateStack(
		middleware.Logging,
		middleware.CORS(corsOptions),
	)
	rootRouter.Handle("/api/", middlewareStack(apiHandler))

	// Create an HTTP server
	s := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           rootRouter,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      300 * time.Second,
		IdleTimeout:       10 * time.Second,
		BaseContext:       func(_ net.Listener) context.Context { return ctx },
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine so it doesn't block
	go func() {
		// Start the server
		logger.Info().Msg(fmt.Sprintf("Server listening on port %d", port))
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("http server exited")
		}
	}()

	// Wait for interrupt signal
	<-stop

	// Shutdown the server
	logger.Info().Msg("Shutting down server...")
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := s.Shutdown(ctxShutdown); err != nil {
		logger.Fatal().Err(err).Msg("server shutdown failed")
	}
	logger.Info().Msg("Server shut down")
}
