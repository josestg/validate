package is_test

import (
	"github.com/josestg/validate/is"
	"testing"
)

func TestInt(t *testing.T) {
	t.Run("int", testIntVariant[int]())
	t.Run("int8", testIntVariant[int8]())
	t.Run("int16", testIntVariant[int16]())
	t.Run("int32", testIntVariant[int32]())
	t.Run("int64", testIntVariant[int64]())

	t.Run("uint", testIntVariant[uint]())
	t.Run("uint8", testIntVariant[uint8]())
	t.Run("uint16", testIntVariant[uint16]())
	t.Run("uint32", testIntVariant[uint32]())
	t.Run("uint64", testIntVariant[uint64]())
}

func testIntVariant[V is.Integers]() func(t *testing.T) {
	return func(t *testing.T) {
		tableTests := []struct {
			name           string
			value          V
			wantError      bool
			wantConstraint string
			validator      is.Validator[V]
		}{
			{
				name:           "it should fail when the int is less than the min",
				value:          1,
				wantError:      true,
				wantConstraint: "integer_min",
				validator:      is.Int[V]().Min(2).Compose(),
			},
			{
				name:           "it should fail when the int is greater than the max",
				value:          3,
				wantError:      true,
				wantConstraint: "integer_max",
				validator:      is.Int[V]().Max(2).Compose(),
			},
			{
				name:           "it should fail when the int is not in choices",
				value:          3,
				wantError:      true,
				wantConstraint: "integer_choose",
				validator:      is.Int[V]().Choose(1, 2).Compose(),
			},
			{
				name:           "it should fail when first validator fails, and the next validator is not called",
				value:          1,
				wantError:      true,
				wantConstraint: "integer_min",
				validator:      is.Int[V]().Min(2).Max(3).Choose(1, 2).Compose(),
			},
			{
				name:           "it should fail when second validator fails, and the next validator is not called",
				value:          4,
				wantError:      true,
				wantConstraint: "integer_max",
				validator:      is.Int[V]().Min(2).Max(3).Choose(1, 2).Compose(),
			},
			{
				name:           "it should fail when third validator fails, and the next validator is not called",
				value:          3,
				wantError:      true,
				wantConstraint: "integer_choose",
				validator:      is.Int[V]().Min(2).Max(3).Choose(1, 2).Compose(),
			},
			{
				name:           "it should pass when all validators pass",
				value:          2,
				wantError:      false,
				wantConstraint: "",
				validator:      is.Int[V]().Min(2).Max(3).Choose(1, 2).Compose(),
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
