package validate_test

import (
	"github.com/josestg/validate"
	"github.com/josestg/validate/is"
	"testing"
)

func TestSchema_Validate(t *testing.T) {
	schema := validate.Schema{
		"name":              validate.Bind("bob", is.String().NotBlank().Len(4, 40).Compose()),
		"email":             validate.Bind("bob@mail", is.String().NotBlank().Email().Compose()),
		"password":          validate.Bind("12345", is.String().NotBlank().Len(6, 20).Compose()),
		"age":               validate.Bind(15, is.Int[int]().Min(18).Max(100).Compose()),
		"favourite_numbers": validate.BindSlice([]int{1, 2, 3, 4, 5}, is.Int[int]().Choose(17, 19).Compose()),
		"height":            validate.Bind(1.5, is.Float[float64]().Min(1.6).Max(3.0).Compose()),
	}

	err := schema.Validate()
	if err == nil {
		t.Error("expected error")
	}

	const expected = `{"age":{"constraint":"integer_min","message":"must be greater than or equal to 18","args":{"min":18,"val":15}},"email":{"constraint":"string_email","message":"must be a valid email address"},"favourite_numbers":{"0":{"constraint":"integer_choose","message":"must be one of [17 19]","args":{"choices":[17,19],"val":1}},"1":{"constraint":"integer_choose","message":"must be one of [17 19]","args":{"choices":[17,19],"val":2}},"2":{"constraint":"integer_choose","message":"must be one of [17 19]","args":{"choices":[17,19],"val":3}},"3":{"constraint":"integer_choose","message":"must be one of [17 19]","args":{"choices":[17,19],"val":4}},"4":{"constraint":"integer_choose","message":"must be one of [17 19]","args":{"choices":[17,19],"val":5}}},"height":{"constraint":"float_min","message":"must be greater than or equal to 1.6","args":{"min":1.6,"val":1.5}},"name":{"constraint":"string_len","message":"must be at least 4 characters","args":{"len":3,"max":40,"min":4}},"password":{"constraint":"string_len","message":"must be at least 6 characters","args":{"len":5,"max":20,"min":6}}}`
	got := err.Error()
	if got != expected {
		t.Errorf("expected %s\n, got %s\n", expected, got)
	}

}
