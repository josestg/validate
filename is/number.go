package is

import "fmt"

// Number is a super set of all numeric types.
type Number interface {
	Integers | Floats
}

// minimum is a validator that checks that value is greater than or equal to min.
func minimum[T Number](constraint string, min T) func(T) error {
	return func(n T) error {
		args := map[string]any{
			"min": min,
			"val": n,
		}

		if n < min {
			return NewError(constraint, fmt.Sprintf("must be greater than or equal to %v", min), args)
		}

		return nil
	}
}

// maximum is a validator that checks that value is less than or equal to max.
func maximum[T Number](constraint string, max T) func(T) error {
	return func(n T) error {
		args := map[string]any{
			"max": max,
			"val": n,
		}

		if max < n {
			return NewError(constraint, fmt.Sprintf("must be less than or equal to %v", max), args)
		}

		return nil
	}
}

// choose is a validator that checks that value is one of choices.
func choose[T Number](constraint string, choices []T) func(T) error {
	return func(n T) error {
		args := map[string]any{
			"choices": choices,
			"val":     n,
		}

		for _, choice := range choices {
			if choice == n {
				return nil
			}
		}

		return NewError(constraint, fmt.Sprintf("must be one of %v", choices), args)
	}
}
