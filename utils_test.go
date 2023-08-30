package decorator_test

import "errors"

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
