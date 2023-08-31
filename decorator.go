package decorator

import (
	"fmt"
	"reflect"
)

type Decorated[M any] struct {
	Call M
}

func MustDecorate[M any](method M, scene Scene, scenes ...Scene) Decorated[M] {
	decorated, err := Decorate(method, scene, scenes...)
	if err != nil {
		panic(err)
	}

	return decorated
}

func Decorate[M any](method M, scene Scene, scenes ...Scene) (Decorated[M], error) {
	var err error
	decorated := Decorated[M]{Call: method}
	decorated.Call, err = getCall(decorated.Call, scene)
	if err != nil {
		return decorated, err
	}

	for _, scene := range scenes {
		if decorated.Call, err = getCall(decorated.Call, scene); err != nil {
			return decorated, err
		}
	}

	return decorated, nil
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
