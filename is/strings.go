package is

import (
	"fmt"
	"regexp"
	"strings"
)

// StringComposer is a composer for string validators
type StringComposer Composer[string]

// String is a constructor for StringComposer.
func String() StringComposer { return Identity[Validator[string]] }

// And composes a string validator with another string validator.
func (f StringComposer) And(next Validator[string]) StringComposer {
	return StringComposer(compose(Composer[string](f), next))
}

// Compose builds a single string validator from the composed validators.
func (f StringComposer) Compose() Validator[string] { return f(nop[string]()) }

// NotBlank validates that the string is not blank.
// consider using NotBlankTrim if you want to trim the string before validating.
// Both NotBlank and NotBlankTrim used "string_not_blank" as constraint name.
func (f StringComposer) NotBlank() StringComposer { return f.And(notBlank()) }

// NotBlankTrim validates that the string is not blank after trimming.
// Both NotBlank and NotBlankTrim() used "string_not_blank" as constraint name.
func (f StringComposer) NotBlankTrim() StringComposer { return f.And(notBlankTrim()) }

// Email validates that the string is a valid email address.
// Email uses "string_email" as constraint name.
func (f StringComposer) Email() StringComposer { return f.And(email()) }

// Len validates that the string length is between min and max.
// If min or max is negative, it is ignored.
// Len uses "string_len" as constraint name.
func (f StringComposer) Len(min int, max int) StringComposer { return f.And(strLen(min, max)) }

func notBlank() Validator[string] {
	return func(s string) error {
		if "" == s {
			return NewError("string_not_blank", "must not be blank", nil)
		}

		return nil
	}
}

func notBlankTrim() Validator[string] {
	return func(s string) error {
		return notBlank().Evaluate(strings.TrimSpace(s))
	}
}

func strLen(min int, max int) Validator[string] {
	return func(s string) error {
		args := map[string]any{
			"min": min,
			"max": max,
			"len": len(s),
		}

		// skip min check if min is negative.
		if min > 0 && min > len(s) {
			return NewError("string_len", fmt.Sprintf("must be at least %d characters", min), args)
		}

		// skip max check if max is negative.
		if max > 0 && max < len(s) {
			return NewError("string_len", fmt.Sprintf("must be at most %d characters", max), args)
		}

		return nil
	}
}

var regexEmail = regexp.MustCompile("^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")

func email() Validator[string] {
	return func(s string) error {
		if !regexEmail.MatchString(s) {
			return NewError("string_email", "must be a valid email address", nil)
		}

		return nil
	}
}
