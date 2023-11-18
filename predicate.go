package funk

import "context"

// Predicate represents a predicate (boolean-valued function) of one argument.
type Predicate[T any] func(context.Context, T) (context.Context, bool, error)

// MustPredicate represents Predicate that doesn't return error.
type MustPredicate[T any] func(context.Context, T) (context.Context, bool)

// PurePredicate represents Predicate that doesn't need context, and even meets the requirements of https://en.wikipedia.org/wiki/Pure_function.
type PurePredicate[T any] func(T) (bool, error)

// PureMustPredicate represents Predicate that doesn't return error and doesn't need context.
type PureMustPredicate[T any] func(T) bool

// Must returns a MustPredicate.
func (p Predicate[T]) Must() MustPredicate[T] {
	return func(ctx context.Context, t T) (context.Context, bool) {
		ctx, v, err := p(ctx, t)
		if err != nil {
			panic(err)
		}
		return ctx, v
	}
}

// Pure returns a PurePredicate.
func (p Predicate[T]) Pure() PurePredicate[T] {
	return func(t T) (bool, error) {
		_, v, err := p(context.Background(), t)
		return v, err
	}
}

// Pure returns a PureMustPredicate.
func (p MustPredicate[T]) Pure() PureMustPredicate[T] {
	return func(t T) bool {
		_, v := p(context.Background(), t)
		return v
	}
}

// Must returns a PureMustPredicate.
func (p PurePredicate[T]) Must() PureMustPredicate[T] {
	return func(t T) bool {
		v, err := p(t)
		if err != nil {
			panic(err)
		}
		return v
	}
}

// And returns a composed Predicate that represents a short-circuiting logical AND of this predicate and another.
func (p Predicate[T]) And(other Predicate[T]) Predicate[T] {
	return func(ctx context.Context, t T) (context.Context, bool, error) {
		ctx, v, err := p(ctx, t)
		if v == false || err != nil {
			return ctx, v, err
		}
		return other(ctx, t)
	}
}

// Or returns a composed Predicate that represents a short-circuiting logical OR of this predicate and another.
func (p Predicate[T]) Or(other Predicate[T]) Predicate[T] {
	return func(ctx context.Context, t T) (context.Context, bool, error) {
		ctx, v, err := p(ctx, t)
		if v == true || err != nil {
			return ctx, v, err
		}
		return other(ctx, t)
	}
}

// Not returns a Predicate that represents the logical negation of this predicate.
func (p Predicate[T]) Not() Predicate[T] {
	return func(ctx context.Context, t T) (context.Context, bool, error) {
		ctx, v, err := p(ctx, t)
		return ctx, !v, err
	}
}

// And returns a composed MustPredicate that represents a short-circuiting logical AND of this predicate and another.
func (p MustPredicate[T]) And(other MustPredicate[T]) MustPredicate[T] {
	return func(ctx context.Context, t T) (context.Context, bool) {
		ctx, v := p(ctx, t)
		if v == false {
			return ctx, v
		}
		return other(ctx, t)
	}
}

// Or returns a composed MustPredicate that represents a short-circuiting logical OR of this predicate and another.
func (p MustPredicate[T]) Or(other MustPredicate[T]) MustPredicate[T] {
	return func(ctx context.Context, t T) (context.Context, bool) {
		ctx, v := p(ctx, t)
		if v == true {
			return ctx, v
		}
		return other(ctx, t)
	}
}

// Not returns a MustPredicate that represents the logical negation of this predicate.
func (p MustPredicate[T]) Not() MustPredicate[T] {
	return func(ctx context.Context, t T) (context.Context, bool) {
		ctx, v := p(ctx, t)
		return ctx, !v
	}
}

// And returns a composed PurePredicate that represents a short-circuiting logical AND of this predicate and another.
func (p PurePredicate[T]) And(other PurePredicate[T]) PurePredicate[T] {
	return func(t T) (bool, error) {
		v, err := p(t)
		if v == false || err != nil {
			return v, err
		}
		return other(t)
	}
}

// Or returns a composed PurePredicate that represents a short-circuiting logical OR of this predicate and another.
func (p PurePredicate[T]) Or(other PurePredicate[T]) PurePredicate[T] {
	return func(t T) (bool, error) {
		v, err := p(t)
		if v == true || err != nil {
			return v, err
		}
		return other(t)
	}
}

// Not returns a PurePredicate that represents the logical negation of this predicate.
func (p PurePredicate[T]) Not() PurePredicate[T] {
	return func(t T) (bool, error) {
		v, err := p(t)
		return !v, err
	}
}

// And returns a composed PureMustPredicate that represents a short-circuiting logical AND of this predicate and another.
func (p PureMustPredicate[T]) And(other PureMustPredicate[T]) PureMustPredicate[T] {
	return func(t T) bool {
		return p(t) && other(t)
	}
}

// Or returns a composed PureMustPredicate that represents a short-circuiting logical OR of this predicate and another.
func (p PureMustPredicate[T]) Or(other PureMustPredicate[T]) PureMustPredicate[T] {
	return func(t T) bool {
		return p(t) || other(t)
	}
}

// Not returns a PureMustPredicate that represents the logical negation of this predicate.
func (p PureMustPredicate[T]) Not() PureMustPredicate[T] {
	return func(t T) bool {
		return !p(t)
	}
}

// BiPredicate represents a predicate (boolean-valued function) of two arguments.
type BiPredicate[T, U any] func(context.Context, T, U) (context.Context, bool, error)

// MustBiPredicate represents BiPredicate that doesn't return error.
type MustBiPredicate[T, U any] func(context.Context, T, U) (context.Context, bool)

// PureBiPredicate represents BiPredicate that doesn't need context, and even meets the requirements of https://en.wikipedia.org/wiki/Pure_function.
type PureBiPredicate[T, U any] func(T, U) (bool, error)

// PureMustBiPredicate represents BiPredicate that doesn't return error and doesn't need context.
type PureMustBiPredicate[T, U any] func(T, U) bool

// Must returns a MustBiPredicate.
func (p BiPredicate[T, U]) Must() MustBiPredicate[T, U] {
	return func(ctx context.Context, t T, u U) (context.Context, bool) {
		ctx, v, err := p(ctx, t, u)
		if err != nil {
			panic(err)
		}
		return ctx, v
	}
}

// Pure returns a PureBiPredicate.
func (p BiPredicate[T, U]) Pure() PureBiPredicate[T, U] {
	return func(t T, u U) (bool, error) {
		_, v, err := p(context.Background(), t, u)
		return v, err
	}
}

// Pure returns a PureMustBiPredicate.
func (p MustBiPredicate[T, U]) Pure() PureMustBiPredicate[T, U] {
	return func(t T, u U) bool {
		_, v := p(context.Background(), t, u)
		return v
	}
}

// Must returns a PureMustBiPredicate.
func (p PureBiPredicate[T, U]) Must() PureMustBiPredicate[T, U] {
	return func(t T, u U) bool {
		v, err := p(t, u)
		if err != nil {
			panic(err)
		}
		return v
	}
}

// And returns a composed predicate that represents a short-circuiting logical AND of this predicate and another.
func (p BiPredicate[T, U]) And(other BiPredicate[T, U]) BiPredicate[T, U] {
	return func(ctx context.Context, t T, u U) (context.Context, bool, error) {
		ctx, v, err := p(ctx, t, u)
		if v == false || err != nil {
			return ctx, v, err
		}
		return other(ctx, t, u)
	}
}

// Or returns a composed predicate that represents a short-circuiting logical OR of this predicate and another.
func (p BiPredicate[T, U]) Or(other BiPredicate[T, U]) BiPredicate[T, U] {
	return func(ctx context.Context, t T, u U) (context.Context, bool, error) {
		ctx, v, err := p(ctx, t, u)
		if v == true || err != nil {
			return ctx, v, err
		}
		return other(ctx, t, u)
	}
}

// Not returns a predicate that represents the logical negation of this predicate.
func (p BiPredicate[T, U]) Not() BiPredicate[T, U] {
	return func(ctx context.Context, t T, u U) (context.Context, bool, error) {
		ctx, v, err := p(ctx, t, u)
		return ctx, !v, err
	}
}

// And returns a composed MustBiPredicate that represents a short-circuiting logical AND of this predicate and another.
func (p MustBiPredicate[T, U]) And(other MustBiPredicate[T, U]) MustBiPredicate[T, U] {
	return func(ctx context.Context, t T, u U) (context.Context, bool) {
		ctx, v := p(ctx, t, u)
		if v == false {
			return ctx, v
		}
		return other(ctx, t, u)
	}
}

// Or returns a composed MustBiPredicate that represents a short-circuiting logical OR of this predicate and another.
func (p MustBiPredicate[T, U]) Or(other MustBiPredicate[T, U]) MustBiPredicate[T, U] {
	return func(ctx context.Context, t T, u U) (context.Context, bool) {
		ctx, v := p(ctx, t, u)
		if v == true {
			return ctx, v
		}
		return other(ctx, t, u)
	}
}

// Not returns a MustBiPredicate that represents the logical negation of this predicate.
func (p MustBiPredicate[T, U]) Not() MustBiPredicate[T, U] {
	return func(ctx context.Context, t T, u U) (context.Context, bool) {
		ctx, v := p(ctx, t, u)
		return ctx, !v
	}
}

// And returns a composed PureBiPredicate that represents a short-circuiting logical AND of this predicate and another.
func (p PureBiPredicate[T, U]) And(other PureBiPredicate[T, U]) PureBiPredicate[T, U] {
	return func(t T, u U) (bool, error) {
		v, err := p(t, u)
		if v == false || err != nil {
			return v, err
		}
		return other(t, u)
	}
}

// Or returns a composed PureBiPredicate that represents a short-circuiting logical OR of this predicate and another.
func (p PureBiPredicate[T, U]) Or(other PureBiPredicate[T, U]) PureBiPredicate[T, U] {
	return func(t T, u U) (bool, error) {
		v, err := p(t, u)
		if v == true || err != nil {
			return v, err
		}
		return other(t, u)
	}
}

// Not returns a PureBiPredicate that represents the logical negation of this predicate.
func (p PureBiPredicate[T, U]) Not() PureBiPredicate[T, U] {
	return func(t T, u U) (bool, error) {
		v, err := p(t, u)
		return !v, err
	}
}

// And returns a composed PureMustBiPredicate that represents a short-circuiting logical AND of this predicate and another.
func (p PureMustBiPredicate[T, U]) And(other PureMustBiPredicate[T, U]) PureMustBiPredicate[T, U] {
	return func(t T, u U) bool {
		return p(t, u) && other(t, u)
	}
}

// Or returns a composed PureMustBiPredicate that represents a short-circuiting logical OR of this predicate and another.
func (p PureMustBiPredicate[T, U]) Or(other PureMustBiPredicate[T, U]) PureMustBiPredicate[T, U] {
	return func(t T, u U) bool {
		return p(t, u) || other(t, u)
	}
}

// Not returns a PureMustBiPredicate that represents the logical negation of this predicate.
func (p PureMustBiPredicate[T, U]) Not() PureMustBiPredicate[T, U] {
	return func(t T, u U) bool {
		return !p(t, u)
	}
}
