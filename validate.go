package validate

import (
	"encoding/json"
	"fmt"
	"github.com/josestg/validate/is"
)

type Errors map[string]error

func (e Errors) Error() string { return stringify(e) }

func stringify(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "stringify: " + err.Error()
	}

	return string(b)
}

type Predicate func() error

// Schema is used to define a validation rules with given fields.
type Schema map[string]Predicate

// Validate validates the given data with the schema.
func (s Schema) Validate() error {
	m := make(Errors)
	for field, predicate := range s {
		if err := predicate(); err != nil {
			switch et := err.(type) {
			default:
				m[field] = err
			case *is.InternalError:
				return et
			}
		}
	}

	if len(m) > 0 {
		return m
	}

	return nil
}

func Bind[T any](data T, validator is.Validator[T]) Predicate {
	return func() error { return validator.Evaluate(data) }
}

func BindSlice[T any](data []T, validator is.Validator[T]) Predicate {
	return func() error {
		errs := make(Errors)
		for index, item := range data {
			if err := validator.Evaluate(item); err != nil {
				errs[fmt.Sprint(index)] = err
			}
		}

		if len(errs) > 0 {
			return errs
		}

		return nil
	}
}
