package inspect

import (
	"reflect"
	"testing"
)

type sample struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestDescribeStruct(t *testing.T) {
	tests := []struct {
		name       string
		in         any
		wantType   string
		wantFields []Field
	}{
		{
			name:     "struct value",
			in:       sample{Name: "Ada", Age: 36},
			wantType: "sample",
			wantFields: []Field{
				{Name: "Name", Type: "string", JSON: "name"},
				{Name: "Age", Type: "int", JSON: "age"},
			},
		},
		{name: "non struct", in: 42, wantType: "int", wantFields: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotFields := DescribeStruct(tt.in)
			if gotType != tt.wantType || !reflect.DeepEqual(gotFields, tt.wantFields) {
				t.Fatalf("DescribeStruct() = (%q, %#v), want (%q, %#v)", gotType, gotFields, tt.wantType, tt.wantFields)
			}
		})
	}
}
