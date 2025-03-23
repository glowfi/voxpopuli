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
	postsTransport := post.NewTransport(services.Post)

	routes := []Route{
		// posts api
		{
			Name:        "PaginatedPost",
			HttpMethod:  GET,
			HttpPath:    "/posts",
			HttpHandler: postsTransport.PostsPaginated,
		},
	}

	return &Server{
		routes:   routes,
		services: services,
	}, nil
}

// HTTPHandler returns the HTTP handler for the server.
func (s *Server) HTTPHandler(ctx context.Context) (http.Handler, error) {
	router := http.NewServeMux()

	if err := HTTPRouter(ctx, router, s.routes); err != nil {
		return nil, err
	}

	return router, nil
}

func HTTPRouter(ctx context.Context, router *http.ServeMux, routes []Route) error {
	for _, r := range routes {
		if r.Name == "" {
			return errors.New("empty route name")
		}

		if r.HttpPath == "" || r.HttpMethod == "" {
			continue
		}

		if r.HttpHandler == nil {
			return fmt.Errorf("nil http handler factory: %s", r.Name)
		}

		switch r.HttpMethod {
		case GET, HEAD, POST, PUT, PATCH, DELETE, CONNECT, OPTIONS, TRACE:
			handlerFunc := r.HttpHandler(ctx, router)
			router.HandleFunc(fmt.Sprintf("%s %s", r.HttpMethod, r.HttpPath), handlerFunc)
		default:
			return fmt.Errorf("invalid http method: %s", r.HttpMethod)
		}
	}
	return nil
}
