package middleware

import (
	"log"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

// getRealIP returns the real IP address of the client.
func getRealIP(r *http.Request) string {
	realIP := r.Header.Get("X-Real-IP")
	if realIP == "" {
		realIP = r.Header.Get("X-Forwarded-For")
	}
	if realIP == "" {
		realIP = r.RemoteAddr
	}
	return realIP
}

// colorize returns a colored string representation of the status code.
func colorize(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "\033[32m" + string(rune(statusCode)) + "\033[0m" // Green for 2xx
	case statusCode >= 300 && statusCode < 400:
		return "\033[36m" + string(rune(statusCode)) + "\033[0m" // Cyan for 3xx
	case statusCode >= 400 && statusCode < 500:
		return "\033[33m" + string(rune(statusCode)) + "\033[0m" // Yellow for 4xx
	case statusCode >= 500:
		return "\033[31m" + string(rune(statusCode)) + "\033[0m" // Red for 5xx
	default:
		return string(rune(statusCode))
	}
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(wrapped, r)

		log.Printf("%s %s %s %s %v\n", colorize(wrapped.statusCode), r.Method, r.URL.Path, getRealIP(r), time.Since(start))
	})
}
