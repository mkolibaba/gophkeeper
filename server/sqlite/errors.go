package sqlite

import "errors"

func asType[E error](err error) (e E, _ bool) {
	return e, errors.As(err, &e)
}
