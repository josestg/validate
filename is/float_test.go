package is_test

import (
	"github.com/josestg/validate/is"
	"testing"
)

func TestFloat(t *testing.T) {
	t.Run("float32", testFloatVariant[float32]())
	t.Run("float64", testFloatVariant[float64]())
}

func testFloatVariant[V is.Floats]() func(t *testing.T) {
	return func(t *testing.T) {
		tableTests := []struct {
			name           string
			value          V
			wantError      bool
			wantConstraint string
			validator      is.Validator[V]
		}{
			{
				name:           "it should fail when the value is less than the min",
				value:          1.0,
				wantError:      true,
				wantConstraint: "float_min",
				validator:      is.Float[V]().Min(2.0).Compose(),
			},
			{
				name:           "it should fail when the value is greater than the max",
				value:          3.0,
				wantError:      true,
				wantConstraint: "float_max",
				validator:      is.Float[V]().Max(2.0).Compose(),
			},
			{
				name:           "it should fail when the value is not in choices",
				value:          3.0,
				wantError:      true,
				wantConstraint: "float_choose",
				validator:      is.Float[V]().Choose(1.0, 2.0).Compose(),
			},
			{
				name:           "it should fail when first validator fails, and the next validator is not called",
				value:          1.4,
				wantError:      true,
				wantConstraint: "float_min",
				validator:      is.Float[V]().Min(2.0).Max(3.0).Choose(1.0, 2.0).Compose(),
			},
			{
				name:           "it should fail when second validator fails, and the next validator is not called",
				value:          4.0,
				wantError:      true,
				wantConstraint: "float_max",
				validator:      is.Float[V]().Min(2.0).Max(3.5).Choose(1, 2).Compose(),
			},
			{
				name:           "it should fail when third validator fails, and the next validator is not called",
				value:          3.0,
				wantError:      true,
				wantConstraint: "float_choose",
				validator:      is.Float[V]().Min(2).Max(3).Choose(1, 2).Compose(),
			},
			{
				name:           "it should pass when all validators pass",
				value:          2.0,
				wantError:      false,
				wantConstraint: "",
				validator:      is.Float[V]().Min(2).Max(3).Choose(1, 2).Compose(),
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
}
