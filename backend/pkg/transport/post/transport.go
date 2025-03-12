package post

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
)

const (
	badRequestStatusCode          = http.StatusBadRequest
	okStatusCode                  = http.StatusOK
	internalServerErrorStatusCode = http.StatusInternalServerError
)

//counterfeiter:generate . PostService
type PostService interface {
	PostsPaginated(ctx context.Context, skip int, limit int) ([]models.Post, error)
}

type Transport struct {
	service PostService
}

type responseError struct {
	Messages []string `json:"errors"`
}

func NewTransport(service PostService) *Transport {
	return &Transport{
		service: service,
	}
}

func writeResponseError(w http.ResponseWriter, statusCode int, errMsgs ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errObj := responseError{Messages: errMsgs}

	if err := json.NewEncoder(w).Encode(errObj); err != nil {
		log.Println("json encode error:", err)
	}
}

func parseIntParam(param string, paramName string) (int, error) {
	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", paramName, err)
	}
	if value < 0 {
		return 0, fmt.Errorf("invalid %s: value must be non-negative", paramName)
	}
	return value, nil
}

func (t *Transport) PostsPaginated(ctx context.Context, h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.Context().Err(); err != nil {
			writeResponseError(w, internalServerErrorStatusCode, "request context error")
			return
		}

		skipStr := r.URL.Query().Get("skip")
		if len(skipStr) == 0 {
			writeResponseError(w, badRequestStatusCode, "add a valid skip number")
			return
		}
		skip, err := parseIntParam(skipStr, "skip")
		if err != nil {
			writeResponseError(w, badRequestStatusCode, err.Error())
			return
		}

		limitStr := r.URL.Query().Get("limit")
		if len(limitStr) == 0 {
			writeResponseError(w, badRequestStatusCode, "add a valid limit")
			return
		}
		limit, err := parseIntParam(limitStr, "limit")
		if err != nil {
			writeResponseError(w, badRequestStatusCode, err.Error())
			return
		}

		posts, err := t.service.PostsPaginated(ctx, skip, limit)
		if err != nil {
			writeResponseError(w, internalServerErrorStatusCode, "failed to fetch posts")
			return
		}

		w.WriteHeader(okStatusCode)
		if err := json.NewEncoder(w).Encode(posts); err != nil {
			log.Println("json encode error while fetching posts:", err)
		}
	})
}
