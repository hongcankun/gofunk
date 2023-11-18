package funk

import "context"

// Func represents a function that accepts one argument and produces a result.
// Note: Golang don't support type parameters in methods of receiver, this severely limits expression ability, such
// as chain calls with generic
type Func[T, R any] func(context.Context, T) (context.Context, R, error)

// MustFunc represents Func that doesn't return error.
type MustFunc[T, R any] func(context.Context, T) (context.Context, R)

// PureFunc represents Func that doesn't need context, and even meets the requirements of https://en.wikipedia.org/wiki/Pure_function.
type PureFunc[T, R any] func(T) (R, error)

// PureMustFunc represents Func that doesn't return error and doesn't need context.
type PureMustFunc[T, R any] func(T) R

// Must returns a MustFunc.
func (f Func[T, R]) Must() MustFunc[T, R] {
	return func(ctx context.Context, t T) (context.Context, R) {
		ctx, v, err := f(ctx, t)
		if err != nil {
			panic(err)
		}
		return ctx, v
	}
}

// Pure returns a PureFunc.
func (f Func[T, R]) Pure() PureFunc[T, R] {
	return func(t T) (R, error) {
		_, v, err := f(context.Background(), t)
		return v, err
	}
}

// Pure returns a PureMustFunc.
func (f MustFunc[T, R]) Pure() PureMustFunc[T, R] {
	return func(t T) R {
		_, v := f(context.Background(), t)
		return v
	}
}

// Must returns a PureMustFunc.
func (f PureFunc[T, R]) Must() PureMustFunc[T, R] {
	return func(t T) R {
		v, err := f(t)
		if err != nil {
			panic(err)
		}
		return v
	}
}

// Unary represents a function on a single operand that produces a result of the same type as its operand.
type Unary[T any] Func[T, T]

// MustUnary represents Unary that doesn't return error.
type MustUnary[T any] MustFunc[T, T]

// PureUnary represents Unary that doesn't need context, and even meets the requirements of https://en.wikipedia.org/wiki/Pure_function.
type PureUnary[T any] PureFunc[T, T]

// PureMustUnary represents Unary that doesn't return error and doesn't need context.
type PureMustUnary[T any] PureMustFunc[T, T]

// Must returns a MustUnary.
func (u Unary[T]) Must() MustUnary[T] {
	return func(ctx context.Context, t T) (context.Context, T) {
		ctx, v, err := u(ctx, t)
		if err != nil {
			panic(err)
		}
		return ctx, v
	}
}

// Pure returns a PureUnary.
func (u Unary[T]) Pure() PureUnary[T] {
	return func(t T) (T, error) {
		_, v, err := u(context.Background(), t)
		return v, err
	}
}

// Pure returns a PureMustUnary.
func (u MustUnary[T]) Pure() PureMustUnary[T] {
	return func(t T) T {
		_, v := u(context.Background(), t)
		return v
	}
}

// Must returns a PureMustUnary.
func (u PureUnary[T]) Must() PureMustUnary[T] {
	return func(t T) T {
		v, err := u(t)
		if err != nil {
			panic(err)
		}
		return v
	}
}

// Then returns a composed Unary that first applies this unary to its input, and then applies the after unary to the result.
func (u Unary[T]) Then(after Unary[T]) Unary[T] {
	return func(ctx context.Context, t T) (context.Context, T, error) {
		ctx, v, err := u(ctx, t)
		if err != nil {
			return ctx, v, err
		}
		return after(ctx, v)
	}
}

// Then returns a composed MustUnary that first applies this unary to its input, and then applies the after unary to the result.
func (u MustUnary[T]) Then(after MustUnary[T]) MustUnary[T] {
	return func(ctx context.Context, t T) (context.Context, T) {
		return after(u(ctx, t))
	}
}

// Then returns a composed PureUnary that first applies this unary to its input, and then applies the after unary to the result.
func (u PureUnary[T]) Then(after PureUnary[T]) PureUnary[T] {
	return func(t T) (T, error) {
		v, err := u(t)
		if err != nil {
			return v, err
		}
		return after(v)
	}
}

// Then returns a composed PureMustUnary that first applies this unary to its input, and then applies the after unary to the result.
func (u PureMustUnary[T]) Then(after PureMustUnary[T]) PureMustUnary[T] {
	return func(t T) T {
		return after(u(t))
	}
}

// BiFunc represents a function that accepts two arguments and produces a result.
// Note: Golang don't support type parameters in methods of receiver, this severely limits expression ability, such
// as chain calls with generic
type BiFunc[T, U, R any] func(context.Context, T, U) (context.Context, R, error)

// MustBiFunc represents BiFunc that doesn't return error.
type MustBiFunc[T, U, R any] func(context.Context, T, U) (context.Context, R)

// PureBiFunc represents BiFunc that doesn't need context, and even meets the requirements of https://en.wikipedia.org/wiki/Pure_function.
type PureBiFunc[T, U, R any] func(T, U) (R, error)

// PureMustBiFunc represents BiFunc that doesn't return error and doesn't need context.
type PureMustBiFunc[T, U, R any] func(T, U) R

// Must returns a MustBiFunc.
func (f BiFunc[T, U, R]) Must() MustBiFunc[T, U, R] {
	return func(ctx context.Context, t T, u U) (context.Context, R) {
		ctx, v, err := f(ctx, t, u)
		if err != nil {
			panic(err)
		}
		return ctx, v
	}
}

// Pure returns a PureBiFunc.
func (f BiFunc[T, U, R]) Pure() PureBiFunc[T, U, R] {
	return func(t T, u U) (R, error) {
		_, v, err := f(context.Background(), t, u)
		return v, err
	}
}

// Pure returns a PureMustBiFunc.
func (f MustBiFunc[T, U, R]) Pure() PureMustBiFunc[T, U, R] {
	return func(t T, u U) R {
		_, v := f(context.Background(), t, u)
		return v
	}
}

// Must returns a PureMustBiFunc.
func (f PureBiFunc[T, U, R]) Must() PureMustBiFunc[T, U, R] {
	return func(t T, u U) R {
		v, err := f(t, u)
		if err != nil {
			panic(err)
		}
		return v
	}
}
