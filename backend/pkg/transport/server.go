package transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/glowfi/voxpopuli/backend/pkg/transport/post"
)

// Supported HTTP methods.
const (
	GET     = "GET"
	HEAD    = "HEAD"
	POST    = "POST"
	PUT     = "PUT"
	PATCH   = "PATCH"
	DELETE  = "DELETE"
	CONNECT = "CONNECT"
	OPTIONS = "OPTIONS"
	TRACE   = "TRACE"
)

// HTTPHandlerFactory creates an HTTP handler function.
type HTTPHandlerFactory func(ctx context.Context, h http.Handler) http.HandlerFunc

// Route represents an HTTP route.
type Route struct {
	Name        string
	HttpMethod  string
	HttpPath    string
	HttpHandler HTTPHandlerFactory
}

// Services represents the services used by the server.
type Services struct {
	Post post.PostService
}

// Server represents the HTTP server.
type Server struct {
	routes   []Route
	services Services
}

// NewServer creates a new server.
func NewServer(services Services) (*Server, error) {
	return &Server{
		services: services,
	}, nil
}

// SetupRoutes sets up the routes for the server.
func (s *Server) SetupRoutes() error {
	postsTransport := post.NewTransport(s.services.Post)

	routes := []Route{
		// posts api
		{
			Name:        "PaginatedPost",
			HttpMethod:  GET,
			HttpPath:    "/posts",
			HttpHandler: postsTransport.PostsPaginated,
		},
	}

	s.routes = routes

	return nil
}

// HTTPHandler returns the HTTP handler for the server.
func (s *Server) HTTPHandler(ctx context.Context) (http.Handler, error) {
	router := http.NewServeMux()

	if err := s.setupHTTPRouter(ctx, router); err != nil {
		return nil, err
	}

	return router, nil
}

func (s *Server) setupHTTPRouter(ctx context.Context, router *http.ServeMux) error {
	for _, r := range s.routes {
		if err := s.validateRoute(r); err != nil {
			return err
		}

		handlerFunc := r.HttpHandler(ctx, router)
		router.HandleFunc(r.HttpPath, handlerFunc)
	}
	return nil
}

func (s *Server) validateRoute(r Route) error {
	if r.Name == "" {
		return errors.New("empty route name")
	}

	if r.HttpPath == "" {
		return errors.New("empty HTTP path")
	}

	if r.HttpMethod == "" {
		return errors.New("empty HTTP method")
	}

	if r.HttpHandler == nil {
		return fmt.Errorf("nil HTTP handler factory: %s", r.Name)
	}

	switch r.HttpMethod {
	case GET, HEAD, POST, PUT, PATCH, DELETE, CONNECT, OPTIONS, TRACE:
		return nil
	default:
		return fmt.Errorf("invalid HTTP method: %s", r.HttpMethod)
	}
}
