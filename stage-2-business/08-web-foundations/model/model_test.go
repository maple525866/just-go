package model

import (
	"reflect"
	"testing"
)

func TestCreateArticleRequestValidationTags(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{name: "title", got: createArticleRequestValidationTag("Title"), want: "required,max=80"},
		{name: "body", got: createArticleRequestValidationTag("Body"), want: "required,min=10"},
		{name: "tags", got: createArticleRequestValidationTag("Tags"), want: "dive,required,max=20"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("validation tag = %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func createArticleRequestValidationTag(field string) string {
	fieldInfo, ok := reflect.TypeOf(CreateArticleRequest{}).FieldByName(field)
	if !ok {
		return ""
	}
	return fieldInfo.Tag.Get("validate")
}
