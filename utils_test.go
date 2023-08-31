package decorator_test

import (
	"errors"
	"fmt"
	"strings"
)

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
