package results

import "errors"

// Unwrap attempts to unwrap a value-error pair.
// If there is an error, it panics. Otherwise, it returns the value.
// Useful for syntactic sugar, but only use this if you do not expect
// any errors.
func Unwrap[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// Try returns a Result given a value and an error.
func Try[T any](value T, err error) *errant[T] {
	return &errant[T]{
		Ok:  value,
		Err: err,
	}
}

// --- Implementation ---

type errant[R any] struct {
	Ok  R
	Err error
}

func (e *errant[R]) Out() (R, error) {
	return e.Ok, e.Err
}

func (e *errant[R]) Unwrap() R {
	if e.Err != nil {
		panic(e)
	}
	return e.Ok
}

func (e *errant[R]) UnwrapOr(value R) R {
	if e.Err != nil {
		return value
	}
	return e.Ok
}

func (e *errant[R]) UnwrapOrElse(onError func() R) R {
	if e.Err != nil {
		return onError()
	}
	return e.Ok
}

func (e *errant[R]) Catch(errorType error, onCatch func(*R, error)) *errant[R] {
	if errors.Is(e.Err, errorType) {
		onCatch(&e.Ok, e.Err)
		e.Err = nil
	}
	return e
}

func (e *errant[R]) CatchTop(onCatch func(*R, error)) *errant[R] {
	onCatch(&e.Ok, e.Err)
	e.Err = errors.Unwrap(e.Err)
	return e
}

func (e *errant[R]) CatchAll(onCatch func(*R, error)) R {
	if e.Err != nil {
		onCatch(&e.Ok, e.Err)
	}
	return e.Ok
}
