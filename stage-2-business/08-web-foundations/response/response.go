package response

import (
	"encoding/json"
	"net/http"

	"just-go/stage-2-business/08-web-foundations/model"
)

// JSON writes a JSON response with a status code.
func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// net/http has no useful recovery path after headers are written.
		return
	}
}

// Error writes the shared REST error envelope.
func Error(w http.ResponseWriter, status int, code, message string, fields []model.FieldError) {
	JSON(w, status, model.ErrorResponse{Code: code, Message: message, Fields: fields})
}
