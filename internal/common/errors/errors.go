package errors

import "errors"

func AsType[E error](err error) (e E, _ bool) {
	return e, errors.As(err, &e)
}
