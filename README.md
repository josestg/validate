# Validate

A programmatic and composable validation library for Go, no more struct tags required. Builder Pattern like API, but without using OOP-style. It uses function composition to create complex validation by combining smaller validation functions (rules).


## Installation

```bash
go get github.com/josestg/validate
```

## Motivation

I was looking for a validation library for Go that not use struct tags, I found ozzo-validation. 
But I am not happy with reflection based-checking. Recently Go released generics, I was thinking maybe I can use it to create a validation library that not rely on reflection.

My goal is to create a validation library that is programmatic and composable by smaller validation functions to create more complex validation functions. 
I came up with using Builder Pattern to compose smaller rules, but I am not happy with the API, it is not very Go idiomatic, it's more like Java rather than Go. 
So I decided to use variadic with option function to apply the rules, for a moment until my code looks like this:

```go
err := validate.Schema{
	"Name": validate.String(data.Name, validate.NotBlank(), validate.MinLength(3), validate.MaxLength(10)),
	"Age":  validate.Int(data.Age, validate.Min(18), validate.Max(100)),
}

if err != nil {
    // handle error
}
```

The code looks much better than the Builder Pattern, but I am not happy with repeating package name `validate` for each rule.
In Builder Pattern, it looks like this:

```go
err := validate.Schema{
	"Name": validate.String().NotBlank().MinLength(3).MaxLength(10).Build().Validate(data.Name),
	"Age":  validate.Int().Min(18).Max(100).Build().Validate(data.Age),
}

if err != nil {
    // handle error
}
```

It just uses the package name once, the trade-off is the struct for Rule Builder uses a lot setter-getter methods, it is not very Go idiomatic.

I was thinking maybe I can use function composition to compose the rules to get the best of both worlds, so I came up with this package. 

Using this package, the code looks like this:

```go
schema := validate.Schema{
    "name":              validate.Bind("bob", is.String().NotBlank().Len(4, 40).Compose()),
    "email":             validate.Bind("bob@mail", is.String().NotBlank().Email().Compose()),
    "password":          validate.Bind("12345", is.String().NotBlank().Len(6, 20).Compose()),
    "age":               validate.Bind(15, is.Int[int]().Min(18).Max(100).Compose()),
    "favourite_numbers": validate.BindSlice([]int{1, 2, 3, 4, 5}, is.Int[int]().Choose(17, 19).Compose()),
    "height":            validate.Bind(1.5, is.Float[float64]().Min(1.6).Max(3.0).Compose()),
}

err := schema.Validate()
if err != nil {
// handle error
}
```

It looks like the Builder Pattern, but it is not using setter-getter methods, it is using function composition to compose the rules.
And it is not using reflection, it is using generics instead.
