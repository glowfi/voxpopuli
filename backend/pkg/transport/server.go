package transport

import (
	"context"
	"fmt"
	"net/http"

	posttransport "github.com/glowfi/voxpopuli/backend/pkg/transport/post"
)

type HTTPHandlerFunc func(ctx context.Context, h http.Handler) http.HandlerFunc

type Route struct {
	Name        string
	HttpMethod  string
	HttpPath    string
	HttpHandler HTTPHandlerFunc
}

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

type Services struct {
	Post posttransport.PostService
}

type Transports struct {
	Post *posttransport.Transport
}

type Server struct {
	routes   []Route
	services Services
}

func NewServer(services Services, transports Transports) (*Server, error) {
	routes := []Route{
		{
			Name:        "PaginatedPost",
			HttpMethod:  GET,
			HttpPath:    "/posts",
			HttpHandler: transports.Post.PostsPaginated,
		},
	}

	return &Server{
		routes:   routes,
		services: services,
	}, nil
}

func HTTPRouter(ctx context.Context, router *http.ServeMux, routes []Route) error {
	for _, r := range routes {
		if r.Name == "" {
			return fmt.Errorf("empty route name")
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
			router.HandleFunc(fmt.Sprintf("%v %v", r.HttpMethod, r.HttpPath), handlerFunc)
		default:
			return fmt.Errorf("invalid http method: %s", r.HttpMethod)
		}
	}
	return nil
}

func (s *Server) HTTPHandler(ctx context.Context) (http.Handler, error) {
	router := http.NewServeMux()

	if err := HTTPRouter(ctx, router, s.routes); err != nil {
		return nil, err
	}

	return router, nil
}
