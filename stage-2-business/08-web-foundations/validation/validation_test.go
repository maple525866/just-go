package validation

import (
	"testing"

	"just-go/stage-2-business/08-web-foundations/model"
)

func TestValidateCreateArticleRequest(t *testing.T) {
	tests := []struct {
		name       string
		request    model.CreateArticleRequest
		wantFields []string
	}{
		{
			name:       "valid request",
			request:    model.CreateArticleRequest{Title: "Routing in Go", Body: "Handlers and middleware make HTTP code testable.", Tags: []string{"http", "go"}},
			wantFields: nil,
		},
		{
			name:       "missing required fields",
			request:    model.CreateArticleRequest{},
			wantFields: []string{"title", "body", "tags"},
		},
		{
			name:       "invalid tag element",
			request:    model.CreateArticleRequest{Title: "Routing", Body: "Long enough body", Tags: []string{""}},
			wantFields: []string{"tags[0]"},
		},
		{
			name:       "blank title",
			request:    model.CreateArticleRequest{Title: "   ", Body: "Long enough body", Tags: []string{"go"}},
			wantFields: []string{"title"},
		},
		{
			name:       "blank body",
			request:    model.CreateArticleRequest{Title: "Routing", Body: "          ", Tags: []string{"go"}},
			wantFields: []string{"body"},
		},
		{
			name:       "blank tag element",
			request:    model.CreateArticleRequest{Title: "Routing", Body: "Long enough body", Tags: []string{"   "}},
			wantFields: []string{"tags[0]"},
		},
	}

	validator := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := validator.ValidateCreateArticle(tt.request)
			if len(fields) != len(tt.wantFields) {
				t.Fatalf("field errors = %+v, want fields %v", fields, tt.wantFields)
			}
			for i, want := range tt.wantFields {
				if fields[i].Field != want {
					t.Fatalf("field[%d] = %q, want %q", i, fields[i].Field, want)
				}
				if fields[i].Rule == "" || fields[i].Message == "" {
					t.Fatalf("field error should include rule and message: %+v", fields[i])
				}
			}
		})
	}
}
