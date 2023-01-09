package is

import "encoding/json"

// Validator is a function that validates a value with a given constraint.
// It returns an error if the value is invalid.
type Validator[T any] func(T) error

// Evaluate is syntactic sugar for calling the validator.
func (f Validator[T]) Evaluate(n T) error { return f(n) }

// mergeValidator merges the given validators into a single validator.
func mergeValidator[T any](validators ...Validator[T]) Validator[T] {
	return func(v T) error {
		for _, validator := range validators {
			if err := validator(v); err != nil {
				return err
			}
		}
		return nil
	}
}

// Composer is a function that composes a validator with another validator.
// It is used to create a chain of validators, or build a validator from
// smaller validators.
type Composer[T any] func(f Validator[T]) Validator[T]

// compose composes a composer with a validator.
// Basically, it composes a chain of validators and wraps them into a single
// validator. In the end, the validator is returned is like a union layer.
func compose[T any](composer Composer[T], next Validator[T]) Composer[T] {
	return func(validator Validator[T]) Validator[T] {
		return mergeValidator(composer(validator), next)
	}
}

// Error is a struct that represents the error that is returned by a validator.
type Error struct {
	// Constraint is the name of the constraint that failed, it can be used
	// to identify the constraint that failed and to provide a custom error
	// translation for it.
	Constraint string `json:"constraint"`

	// Message is the error message that is returned by the validator.
	// This message can be used as fallback if no custom error translation
	// is provided.
	Message string `json:"message"`

	// Args is a map of arguments that can be used to provide additional
	// information about the error.
	// For example, the "min" constraint can return the minimum value that
	// was expected.
	Args map[string]any `json:"args,omitempty"`
}

// NewError creates a new error with the given constraint, message and args.
func NewError(constraint string, message string, args map[string]any) *Error {
	return &Error{Constraint: constraint, Message: message, Args: args}
}

// Error returns the error message.
func (e *Error) Error() string { return stringify(e) }

// InternalError is a struct that represents an service error.
// It's used to tell validator this kind of error is not expected to happen,
// And it should be returned instead of the validation error.
//
// When validator found an service error, it will immediately to the caller
// And validation chain will be stopped.
type InternalError struct {
	Key string
	Err error
}

// NewInternalError creates a new service error with the given key And error.
func NewInternalError(key string, err error) *InternalError {
	return &InternalError{Key: key, Err: err}
}

func (e *InternalError) Error() string { return e.Key + ": " + e.Err.Error() }

// stringify converts the given value to a JSON string.
func stringify(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}

	return string(b)
}

// Identity is a function that returns the given value.
func Identity[F any](f F) F { return f }

// nop is a function that returns a validator that always returns nil.
// It is used as terminal function for the composer chain.
func nop[T any]() Validator[T] { return func(_ T) error { return nil } }
