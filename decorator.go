package decorator

import (
	"fmt"
	"reflect"
)

func MustGetCall[D, M any](decor Scene[D], method M) M {
	call, err := GetCall(decor, method)
	if err != nil {
		panic(err)
	}

	return call
}

func GetCall[D, M any](decor Scene[D], method M) (M, error) {
	var zero M
	t := reflect.TypeOf(method)
	if t.Kind() != reflect.Func {
		return zero, fmt.Errorf("wrong method argument %s: %w", t, ErrNotAFunc)
	}

	registryElement, ok := decor.GetCall(t.String())
	if !ok {
		return method, nil
	}

	decoratedFn, ok := registryElement.(func(D, M) M)

	if !ok {
		return zero, fmt.Errorf(
			"expected %s, got %s: %w",
			reflect.TypeOf(func(D, M) M { return zero }),
			reflect.TypeOf(registryElement),
			ErrWrongDecoratedType)
	}

	return decoratedFn(decor.D(), method), nil
}
