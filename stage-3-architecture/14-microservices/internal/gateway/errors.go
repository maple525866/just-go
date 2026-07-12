package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeGatewayError(w http.ResponseWriter, err error) {
	code := status.Code(err)
	switch {
	case errors.Is(err, context.DeadlineExceeded), code == codes.DeadlineExceeded:
		writeJSONError(w, http.StatusGatewayTimeout, "downstream request timed out")
	case errors.Is(err, context.Canceled), code == codes.Canceled:
		writeJSONError(w, http.StatusGatewayTimeout, "downstream request canceled")
	case code == codes.InvalidArgument:
		writeJSONError(w, http.StatusBadRequest, "invalid downstream request")
	case code == codes.NotFound:
		writeJSONError(w, http.StatusNotFound, "resource not found")
	case code == codes.Unavailable:
		writeJSONError(w, http.StatusServiceUnavailable, "service unavailable")
	default:
		writeJSONError(w, http.StatusInternalServerError, "request failed")
	}
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: message})
}
