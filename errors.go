package decorator

import "errors"

var (
	ErrNotAFunc           = errors.New("must be a method or a function")
	ErrWrongDecoratedType = errors.New("requested function call has wrong decoration type registered")
	ErrDuplicate          = errors.New("decoration was already added before")
)
