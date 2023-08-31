package decorator

import (
	"fmt"
	"reflect"
)

// Scene constructor option.
type Option[D any] func(map[string]any) error

// Scene constructor option to provide decoration mechanism for function or method call.
func SceneDecor[M, D any](fn M, decoration func(d D, originalCall M) M) Option[D] {
	return func(opts map[string]any) error {
		t := reflect.TypeOf(fn)
		if t.Kind() != reflect.Func {
			return fmt.Errorf("failed to decorate %s: %w", t, ErrNotAFunc)
		}

		fnName := t.String()
		if _, ok := opts[fnName]; ok {
			return fmt.Errorf("failed to decorate %s: %w", t, ErrDuplicate)
		}

		opts[fnName] = func(v any, fn M) M { return decoration(v.(D), fn) }

		return nil
	}
}
