# decorator

This package provides variation of Decorator pattern.

### To install tinysl:
`go get -u github.com/andriiyaremenko/decorator`

### How to use:
```go
type someOtherService struct{ s string }

func (t *someOtherService) SomeMethod(prefix string) string {
	return prefix + t.s
}

func (t *someOtherService) S() string {
	return t.s
}

type someService struct{ s string }

func (t *someService) SomeMethod(prefix string, fail bool) (string, error) {
	if !fail {
		return prefix + t.s, nil
	}

	return "", errors.New("some error")
}

type validate struct{}

func (t *validate) Invalid(s string) error {
	if s == "fail" {
		return errors.New("validation failed")
	}

	return nil
}

// ... other file

var (
	service                      = &someService{s: "test"}
	anotherService               = &someOtherService{s: "fail"}
	anotherServiceComposedMethod func(*someOtherService, string) (string, error)
	validateDecorator            decorator.Scene[*validate]
)

validateDecorator, err := decorator.NewScene(
	&validate{},
	decorator.Decorate((*someService).SomeMethod,
		func(
			d *validate, fn func(*someService, string, bool) (string, error),
		) func(*someService, string, bool) (string, error) {
			return func(s *someService, prefix string, fail bool) (string, error) {
				if err := d.Invalid(prefix); err != nil {
					return "", err
				}

				return fn(s, prefix, fail)
			}
		}),
	decorator.Decorate(anotherServiceComposedMethod,
		func(
			d *validate, fn func(*someOtherService, string) (string, error),
		) func(*someOtherService, string) (string, error) {
			return func(s *someOtherService, prefix string) (string, error) {
				if err := d.Invalid(s.S()); err != nil {
					return "", err
				}

				return s.SomeMethod(prefix), nil
			}
		}),
)

result, err := decorator.MustGetCall(validateDecorator, (*someService).SomeMethod)(service, "some ", false)
anotherResult, err := decorator.MustGetCall(validateDecorator, anotherServiceComposedMethod)(anotherService, "some ")

// check errors, use results...
```

### Types
 * `decorator.Scene` - type representing decorator consisted of decorating type and registry of decorated methods and functions.
 * `decorator.Option` - options type to use with `decorator.NewScene`.

### Functions:
 * `decorator.GetCall`
 * `decorator.MustGetCall`
 * `decorator.NewScene`

### Options for new Scene:
 * `decorator.Decorate`
