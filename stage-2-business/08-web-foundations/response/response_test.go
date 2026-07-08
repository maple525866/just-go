package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"just-go/stage-2-business/08-web-foundations/model"
)

func TestJSONWritesStatusHeaderAndBody(t *testing.T) {
	recorder := httptest.NewRecorder()
	payload := map[string]string{"status": "ok"}

	JSON(recorder, http.StatusAccepted, payload)

	if recorder.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusAccepted)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content type = %q, want application/json", got)
	}
	var body map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("body status = %q, want ok", body["status"])
	}
}

func TestErrorWritesEnvelope(t *testing.T) {
	recorder := httptest.NewRecorder()
	fields := []model.FieldError{{Field: "title", Rule: "required", Message: "title is required"}}

	Error(recorder, http.StatusUnprocessableEntity, "validation_failed", "request validation failed", fields)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnprocessableEntity)
	}
	var body model.ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if body.Code != "validation_failed" || len(body.Fields) != 1 {
		t.Fatalf("unexpected error body: %+v", body)
	}
}
