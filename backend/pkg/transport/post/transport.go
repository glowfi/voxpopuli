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

//counterfeiter:generate . PostService
type PostService interface {
	PostsPaginated(ctx context.Context, skip, limit int) ([]models.PostPaginated, error)
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

func (t *Transport) PostsPaginated(w http.ResponseWriter, r *http.Request) {
	if err := r.Context().Err(); err != nil {
		writeResponseError(w, http.StatusInternalServerError, "request context error")
		return
	}

	skipStr := r.URL.Query().Get("skip")
	if len(skipStr) == 0 {
		writeResponseError(w, http.StatusBadRequest, "add a valid skip number")
		return
	}
	skip, err := parseIntParam(skipStr, "skip")
	if err != nil {
		writeResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	limitStr := r.URL.Query().Get("limit")
	if len(limitStr) == 0 {
		writeResponseError(w, http.StatusBadRequest, "add a valid limit")
		return
	}
	limit, err := parseIntParam(limitStr, "limit")
	if err != nil {
		writeResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	posts, err := t.service.PostsPaginated(r.Context(), skip, limit)
	if err != nil {
		writeResponseError(w, http.StatusInternalServerError, "failed to fetch posts")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Println("json encode error while fetching posts:", err)
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
