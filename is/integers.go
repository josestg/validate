package is

// Integers is a type alias for the set of all integer types.
type Integers interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// IntComposer is composer for integer validators.
type IntComposer[T Integers] Composer[T]

// Int is constructor for integer validators.
func Int[T Integers]() IntComposer[T] {
	return Identity[Validator[T]]
}

// add composes a validator with the composer.
func (f IntComposer[T]) and(next Validator[T]) IntComposer[T] {
	return IntComposer[T](compose(Composer[T](f), next))
}

// Compose composes all validators into one validator.
func (f IntComposer[T]) Compose() Validator[T] { return f(nop[T]()) }

// Min checks that value is greater than or equal to min.
// Min is "integer_min" as constraint name
func (f IntComposer[T]) Min(min T) IntComposer[T] { return f.and(minimum("integer_min", min)) }

// Max checks that value is less than or equal to max.
// Max is "integer_max" as constraint name.
func (f IntComposer[T]) Max(max T) IntComposer[T] { return f.and(maximum("integer_max", max)) }

// Choose checks that value is one of choices.
// Choose is "integer_choose" as constraint name.
func (f IntComposer[T]) Choose(choices ...T) IntComposer[T] {
	return f.and(choose("integer_choose", choices))
}
