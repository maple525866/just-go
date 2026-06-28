package apperr

import (
	"errors"
	"strings"
	"testing"
)

func TestFindUserAndErrorsIs(t *testing.T) {
	tests := []struct {
		name       string
		user       string
		wantValue  string
		wantErrIs  bool
		wantErrNil bool
	}{
		{name: "known user", user: "Ada", wantValue: "Ada Lovelace", wantErrNil: true},
		{name: "missing user", user: "Zoe", wantErrIs: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindUser(tt.user)
			if got != tt.wantValue {
				t.Fatalf("FindUser() value = %q, want %q", got, tt.wantValue)
			}
			if (err == nil) != tt.wantErrNil {
				t.Fatalf("FindUser() err nil = %t, want %t", err == nil, tt.wantErrNil)
			}
			if IsUserNotFound(err) != tt.wantErrIs {
				t.Fatalf("IsUserNotFound() = %t, want %t", IsUserNotFound(err), tt.wantErrIs)
			}
			if errors.Is(err, ErrUserNotFound) != tt.wantErrIs {
				t.Fatalf("errors.Is() = %t, want %t", errors.Is(err, ErrUserNotFound), tt.wantErrIs)
			}
		})
	}
}

func TestExtractQueryError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		wantFound bool
		wantUser  string
	}{
		{
			name: "wrapped query error",
			err: func() error {
				_, err := FindUser("Zoe")
				return err
			}(),
			wantFound: true,
			wantUser:  "Zoe",
		},
		{name: "plain sentinel", err: ErrUserNotFound, wantFound: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := ExtractQueryError(tt.err)
			if found != tt.wantFound {
				t.Fatalf("ExtractQueryError() found = %t, want %t", found, tt.wantFound)
			}
			if found && got.User != tt.wantUser {
				t.Fatalf("ExtractQueryError() user = %q, want %q", got.User, tt.wantUser)
			}
		})
	}
}

func TestSummary(t *testing.T) {
	tests := []struct {
		name     string
		contains []string
	}{
		{name: "mentions wrapping and checks", contains: []string{"%w", "errors.Is", "errors.As"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Summary()
			for _, part := range tt.contains {
				if !strings.Contains(got, part) {
					t.Fatalf("Summary() = %q, want it to contain %q", got, part)
				}
			}
		})
	}
}
