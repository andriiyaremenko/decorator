package decorator

import (
	"fmt"
	"reflect"
)

type Option[D any] func(map[string]any) error

func Decorate[M, D any](fn M, decoration func(d D, originalCall M) M) Option[D] {
	return func(opts map[string]any) error {
		t := reflect.TypeOf(fn)
		if t.Kind() != reflect.Func {
			return fmt.Errorf("failed to decorate %s: %w", t, ErrNotAFunc)
		}

		fnName := t.String()
		if _, ok := opts[fnName]; ok {
			return fmt.Errorf("failed to decorate %s: %w", t, ErrDuplicate)
		}

		opts[fnName] = decoration

		return nil
	}
}
