package codec

import (
	"reflect"
	"strings"
	"testing"
)

func TestJSONRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		lesson   Lesson
		contains string
	}{
		{name: "json", lesson: Lesson{ID: 6, Title: "stdlib", Done: true}, contains: `"title":"stdlib"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, decoded, err := JSONRoundTrip(tt.lesson)
			if err != nil {
				t.Fatalf("JSONRoundTrip() unexpected error: %v", err)
			}
			if !strings.Contains(encoded, tt.contains) || !reflect.DeepEqual(decoded, tt.lesson) {
				t.Fatalf("JSONRoundTrip() = (%q, %#v), want encoded containing %q and decoded %#v", encoded, decoded, tt.contains, tt.lesson)
			}
		})
	}
}

func TestXMLRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		lesson   Lesson
		contains string
	}{
		{name: "xml", lesson: Lesson{ID: 6, Title: "stdlib", Done: true}, contains: `<lesson id="6">`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, decoded, err := XMLRoundTrip(tt.lesson)
			if err != nil {
				t.Fatalf("XMLRoundTrip() unexpected error: %v", err)
			}
			wantDecoded := tt.lesson
			wantDecoded.XMLName.Local = "lesson"
			if !strings.Contains(encoded, tt.contains) || !reflect.DeepEqual(decoded, wantDecoded) {
				t.Fatalf("XMLRoundTrip() = (%q, %#v), want encoded containing %q and decoded %#v", encoded, decoded, tt.contains, wantDecoded)
			}
		})
	}
}
