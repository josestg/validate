package is_test

import (
	"github.com/josestg/validate/is"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	tableTests := []struct {
		name           string
		value          string
		wantError      bool
		wantConstraint string
		validator      is.Validator[string]
	}{
		{
			name:           "it should fail when the string is blank",
			value:          "",
			wantError:      true,
			wantConstraint: "string_not_blank",
			validator:      is.String().NotBlank().Compose(),
		},
		{
			name:           "it should fail when the string is blank after trimming",
			value:          " ",
			wantError:      true,
			wantConstraint: "string_not_blank",
			validator:      is.String().NotBlankTrim().Compose(),
		},
		{
			name:           "it should fail when the string is not an email",
			value:          "bob",
			wantError:      true,
			wantConstraint: "string_email",
			validator:      is.String().Email().Compose(),
		},
		{
			name:           "it should fail when the string is too short",
			value:          "bob",
			wantError:      true,
			wantConstraint: "string_len",
			validator:      is.String().Len(4, 10).Compose(),
		},
		{
			name:           "it should fail when the string is too long",
			value:          strings.Repeat("a", 11),
			wantError:      true,
			wantConstraint: "string_len",
			validator:      is.String().Len(4, 10).Compose(),
		},
		{
			name:           "it should fail when first validator fails, and the next validator is not called",
			value:          "",
			wantError:      true,
			wantConstraint: "string_not_blank",
			validator:      is.String().NotBlank().Len(4, 10).Email().Compose(),
		},
		{
			name:           "it should fail when second validator fails, and the next validator is not called",
			value:          "bob",
			wantError:      true,
			wantConstraint: "string_len",
			validator:      is.String().NotBlank().Len(4, 10).Email().Compose(),
		},
		{
			name:           "it should fail when third validator fails, and the next validator is not called",
			value:          "bob alexander",
			wantError:      true,
			wantConstraint: "string_email",
			validator:      is.String().NotBlank().Len(4, 20).Email().Compose(),
		},
		{
			name:           "it should pass when all validators pass",
			value:          "bob@mail.com",
			wantError:      false,
			wantConstraint: "",
			validator:      is.String().NotBlank().Len(4, 20).Email().Compose(),
		},
	}

	for _, tt := range tableTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validator(tt.value)
			if tt.wantError {
				verifyValidateError(t, err, tt.wantConstraint)
			} else {
				if nil != err {
					t.Fatal(err)
				}
			}
		})
	}

}

func verifyValidateError(t *testing.T, err error, constraint string) {
	if nil == err {
		t.Fatal("expected error, got nil")
	}

	e, ok := err.(*is.Error)
	if !ok {
		t.Fatalf("expected validate.Error, got %T", err)
	}

	t.Log(e.Error())
	if e.Constraint != constraint {
		t.Fatalf("expected constraint %q, got %q", constraint, e.Constraint)
	}
}
