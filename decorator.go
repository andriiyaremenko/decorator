package decorator

import (
	"fmt"
	"reflect"
)

func MustDecorate[M any](method M, scene Scene, scenes ...Scene) M {
	call, err := Decorate(method, scene, scenes...)
	if err != nil {
		panic(err)
	}

	return call
}

func Decorate[M any](method M, scene Scene, scenes ...Scene) (M, error) {
	method, err := getCall(method, scene)
	if err != nil {
		return method, err
	}

	for _, scene := range scenes {
		if method, err = getCall(method, scene); err != nil {
			return method, err
		}
	}

	return method, nil
}

func getCall[M any](method M, scene Scene) (M, error) {
	var zero M
	t := reflect.TypeOf(method)
	if t.Kind() != reflect.Func {
		return zero, fmt.Errorf("wrong method argument %s: %w", t, ErrNotAFunc)
	}

	registryElement, ok := scene.GetCall(t.String())
	if !ok {
		return method, nil
	}

	decoratedFn, ok := registryElement.(func(any, M) M)

	if !ok {
		return zero, fmt.Errorf(
			"expected %s, got %s: %w",
			reflect.TypeOf(func(any, M) M { return zero }),
			reflect.TypeOf(registryElement),
			ErrWrongDecoratedType)
	}

	return decoratedFn(scene.D(), method), nil
}
