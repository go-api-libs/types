package types_test

import (
	"strings"
	"testing"

	"github.com/go-api-libs/types"
	"github.com/go-json-experiment/json"
)

type myStruct struct {
	Email types.Email `json:"email,omitempty"`
}

type myStruct2 struct {
	Email *types.Email `json:"email,omitempty"`
}

func TestEmail(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name      string
		json      string
		want      string
		wantError string
	}{
		{"empty", `{"email":""}`, "", ""},
		{"fully empty", `{}`, "", ""},
		{"valid", `{"email":"max@example.com"}`, "max@example.com", ""},
		{
			"invalid", `{"email":"foo"}`, "",
			`json: cannot unmarshal Go value of type types.Email: invalid email: "foo"`,
		},
		{"invalid", `{"email":"foo`, "", `json: cannot unmarshal Go value of type types.Email: unexpected EOF`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			t.Run("no pointer", func(t *testing.T) {
				res := &myStruct{}
				err := json.Unmarshal([]byte(tc.json), res)
				if tc.wantError != "" {
					if err == nil {
						t.Fatal("expected error")
					} else if got := err.Error(); strings.Replace(got,
						// sometimes the error message is "unable to" and sometimes "cannot"
						"unable to", "cannot", 1) != tc.wantError {
						t.Fatalf("got: %q, want: %q", got, tc.wantError)
					}
					return
				}
				if err != nil {
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
				res := &myStruct2{}
				err := json.Unmarshal([]byte(tc.json), res)
				if tc.wantError != "" {
					if err == nil {
						t.Fatal("expected error")
					} else if got := err.Error(); strings.Replace(got,
						// sometimes the error message is "unable to" and sometimes "cannot"
						"unable to", "cannot", 1) != tc.wantError {
						t.Fatalf("got: %q, want: %v", got, tc.wantError)
					}
					return
				}
				if err != nil {
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
