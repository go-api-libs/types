package types_test

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"reflect"
	"testing"

	"github.com/go-api-libs/types"
)

type wrapEmail struct {
	Email types.Email `json:"email,omitempty"`
}

type wrapPointerToEmail struct {
	Email *types.Email `json:"email,omitempty"`
}

func TestEmail(t *testing.T) {
	for _, tc := range []struct {
		name string
		json string
		want string
	}{
		{"empty", `{"email":""}`, ""},
		{"fully empty", `{}`, ""},
		{"valid", `{"email":"max@example.com"}`, "max@example.com"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			t.Run("no pointer", func(t *testing.T) {
				res := &wrapEmail{}
				if err := json.Unmarshal([]byte(tc.json), res); err != nil {
					t.Fatal(err)
				}

				if res.Email != types.Email(tc.want) {
					t.Fatalf("got: %v, want: %v", res.Email, tc.want)
				}

				if tc.want != "" {
					if err := res.Email.Validate(); err != nil {
						t.Fatal(err)
					}
				}
			})

			t.Run("pointer", func(t *testing.T) {
				res := &wrapPointerToEmail{}
				if err := json.Unmarshal([]byte(tc.json), res); err != nil {
					t.Fatal(err)
				}

				if (res.Email == nil && tc.want != "") ||
					(res.Email != nil && *res.Email != types.Email(tc.want)) {
					t.Fatalf("got: %#v, want: %q", res.Email, tc.want)
				}

				if tc.want != "" {
					if err := res.Email.Validate(); err != nil {
						t.Fatal(err, res.Email == nil)
					}
				}
			})
		})
	}
}

func TestEmail_Errors(t *testing.T) {
	emailType := reflect.TypeFor[types.Email]()

	t.Run("invalid email", func(t *testing.T) {
		body := []byte(`{"email":"foo"}`)
		for _, out := range []any{&wrapEmail{}, &wrapPointerToEmail{}} {
			semErr := &json.SemanticError{}
			if err := json.Unmarshal(body, out); err == nil {
				t.Fatal("expected error")
			} else if !errors.As(err, &semErr) {
				t.Fatalf("wanted: %T, got: %T", semErr, err)
			} else if semErr.GoType != emailType {
				t.Fatalf("wanted: %T, got: %T", emailType, semErr.GoType)
			}
		}
	})

	t.Run("unexpected EOF", func(t *testing.T) {
		body := []byte(`{"email":"foo`)
		for _, out := range []any{&wrapEmail{}, &wrapPointerToEmail{}} {
			synErr := &jsontext.SyntacticError{}
			if err := json.Unmarshal(body, out); err == nil {
				t.Fatal("expected error")
			} else if !errors.As(err, &synErr) {
				t.Fatalf("wanted: %T, got: %T", synErr, err)
			} else if want := `unexpected EOF`; synErr.Err.Error() != want {
				t.Fatalf("wanted: %s, got: %s", want, synErr.Err)
			}
		}
	})
}
