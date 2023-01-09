package is

// Floats is a generic interface for float32 and float64.
type Floats interface {
	~float32 | ~float64
}

// FloatComposer is composer for float validators.
type FloatComposer[T Floats] Composer[T]

// Float is constructor for float composer.
func Float[T Floats]() FloatComposer[T] { return Identity[Validator[T]] }

// add composes validator with next validator.
func (f FloatComposer[T]) and(next Validator[T]) FloatComposer[T] {
	return FloatComposer[T](compose(Composer[T](f), next))
}

// Compose composes all validators into one validator.
func (f FloatComposer[T]) Compose() Validator[T] { return f(nop[T]()) }

// Min checks that value is greater than or equal to min.
// Min is "float_min" as constraint name.
func (f FloatComposer[T]) Min(min T) FloatComposer[T] { return f.and(minimum("float_min", min)) }

// Max checks that value is less than or equal to max.
// Max is "float_max" as constraint name.
func (f FloatComposer[T]) Max(max T) FloatComposer[T] { return f.and(maximum("float_max", max)) }

// Choose checks that value is one of choices.
// Choose is "float_choose" as constraint name.
func (f FloatComposer[T]) Choose(choices ...T) FloatComposer[T] {
	return f.and(choose("float_choose", choices))
}
