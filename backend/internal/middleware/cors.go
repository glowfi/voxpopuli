package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// CORSOptions represents the options for the CORS middleware.
type CORSOptions struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
	ExposeHeaders  []string
	MaxAge         int
}

// DefaultCORSOptions returns the default CORS options.
func DefaultCORSOptions() CORSOptions {
	return CORSOptions{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:  []string{"Content-Length"},
		MaxAge:         3600,
	}
}

// CORS returns a CORS middleware with the given options.
func CORS(options CORSOptions) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			for _, allowedOrigin := range options.AllowedOrigins {
				if origin == allowedOrigin || allowedOrigin == "*" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}

			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(options.AllowedMethods, ","))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(options.AllowedHeaders, ","))
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(options.ExposeHeaders, ","))
				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(options.MaxAge))
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
