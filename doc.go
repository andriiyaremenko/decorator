/*
This package provides variation of Decorator pattern.

To install decorator:

	go get -u github.com/andriiyaremenko/decorator

How to use:

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

	type log struct{ s strings.Builder }

	func (l *log) GetLog() string {
		return l.s.String()
	}

	func (l *log) LogSomeServiceSomeMethod(
		fn func(*someService, string, bool) (string, error),
	) func(*someService, string, bool) (string, error) {
		return func(s *someService, arg1 string, arg2 bool) (string, error) {
			l.s.WriteString(fmt.Sprintf("got: %q, %t\t'", arg1, arg2))
			result, err := fn(s, arg1, arg2)
			if err != nil {
				l.s.WriteString(fmt.Sprintf("resulted in error: %s\n", err))
				return result, err
			}

			l.s.WriteString(fmt.Sprintf("resulted in %q\n", result))
			return result, err
		}
	}

	func (l *log) LogSomeOtherServiceSomeMethod(
		fn func(*someOtherService, string) string,
	) func(*someOtherService, string) string {
		return func(s *someOtherService, arg string) string {
			l.s.WriteString(fmt.Sprintf("got: %q\t'", arg))

			result := fn(s, arg)

			l.s.WriteString(fmt.Sprintf("resulted in %q\n", result))

			return result
		}
	}

	// ... other file

	var (
		service                      = &someService{s: "test"}
		anotherService               = &someOtherService{s: "fail"}
		anotherServiceComposedMethod func(*someOtherService, string) (string, error)
		validateDecorator            decorator.Scene
		logDecorator                 decorator.Scene
	)

	validateDecorator, err := decorator.NewScene(
		&validate{},
		decorator.SceneDecor((*someService).SomeMethod,
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
		decorator.SceneDecor(anotherServiceComposedMethod,
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

	logDecorator, err = decorator.NewScene(
		&log{s: strings.Builder{}},
		decorator.SceneDecor((*someService).SomeMethod, (*log).LogSomeServiceSomeMethod),
		decorator.SceneDecor((*someOtherService).SomeMethod, (*log).LogSomeOtherServiceSomeMethod),
	)

	result, err := decorator.MustDecorate(
		(*someService).SomeMethod, validateDecorator, logDecorator
	)(service, "some ", false)
	anotherResult, err := decorator.MustDecorate(
		anotherServiceComposedMethod, validateDecorator, logDecorator
	)(anotherService, "some ")

	// check errors, use results...

Types
  - decorator.Scene - type representing decorator consisted of decorating type and registry of decorated methods and functions.
  - decorator.Option - options type to use with decorator.NewScene.

Functions:
  - decorator.Decorate
  - decorator.MustDecorate
  - decorator.NewScene

Options for new Scene:
  - decorator.SceneDecor
*/
package decorator
