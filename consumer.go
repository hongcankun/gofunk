package funk

import "context"

// Consumer represents an operation that accepts a single input argument and returns no result.
type Consumer[T any] func(context.Context, T) (context.Context, error)

// MustConsumer represents Consumer that doesn't return error.
type MustConsumer[T any] func(context.Context, T) context.Context

// PureConsumer represents Consumer that doesn't need context, and even meets the requirements of https://en.wikipedia.org/wiki/Pure_function.
type PureConsumer[T any] func(T) error

// PureMustConsumer represents Consumer that doesn't return error and doesn't need context.
type PureMustConsumer[T any] func(T)

// Must returns a MustConsumer.
func (c Consumer[T]) Must() MustConsumer[T] {
	return func(ctx context.Context, t T) context.Context {
		ctx, err := c(context.Background(), t)
		if err != nil {
			panic(err)
		}
		return ctx
	}
}

// Pure returns a PureConsumer.
func (c Consumer[T]) Pure() PureConsumer[T] {
	return func(t T) error {
		_, err := c(context.Background(), t)
		return err
	}
}

// Pure returns a PureMustConsumer.
func (c MustConsumer[T]) Pure() PureMustConsumer[T] {
	return func(t T) {
		_ = c(context.Background(), t)
		return
	}
}

// Must returns a PureMustConsumer.
func (c PureConsumer[T]) Must() PureMustConsumer[T] {
	return func(t T) {
		err := c(t)
		if err != nil {
			panic(err)
		}
		return
	}
}

// Then returns a composed Consumer that performs, in sequence, this operation followed by the after operation.
func (c Consumer[T]) Then(after Consumer[T]) Consumer[T] {
	return func(ctx context.Context, t T) (context.Context, error) {
		ctx, err := c(ctx, t)
		if err != nil {
			return ctx, err
		}
		return after(ctx, t)
	}
}

// Then returns a composed MustConsumer that performs, in sequence, this operation followed by the after operation.
func (c MustConsumer[T]) Then(after MustConsumer[T]) MustConsumer[T] {
	return func(ctx context.Context, t T) context.Context {
		return after(c(ctx, t), t)
	}
}

// Then returns a composed PureConsumer that performs, in sequence, this operation followed by the after operation.
func (c PureConsumer[T]) Then(after PureConsumer[T]) PureConsumer[T] {
	return func(t T) error {
		err := c(t)
		if err != nil {
			return err
		}
		return after(t)
	}
}

// Then returns a composed PureMustConsumer that performs, in sequence, this operation followed by the after operation.
func (c PureMustConsumer[T]) Then(after PureMustConsumer[T]) PureMustConsumer[T] {
	return func(t T) {
		c(t)
		after(t)
	}
}

// BiConsumer represents a function that accepts two arguments and produces a result.
type BiConsumer[T, U any] func(context.Context, T, U) (context.Context, error)

// MustBiConsumer represents BiConsumer that doesn't return error.
type MustBiConsumer[T, U any] func(context.Context, T, U) context.Context

// PureBiConsumer represents BiConsumer that doesn't need context, and even meets the requirements of https://en.wikipedia.org/wiki/Pure_function.
type PureBiConsumer[T, U any] func(T, U) error

// PureMustBiConsumer represents BiConsumer that doesn't return error and doesn't need context.
type PureMustBiConsumer[T, U any] func(T, U)

// Must return a MustBiConsumer.
func (c BiConsumer[T, U]) Must() MustBiConsumer[T, U] {
	return func(ctx context.Context, t T, u U) context.Context {
		ctx, err := c(ctx, t, u)
		if err != nil {
			panic(err)
		}
		return ctx
	}
}

// Pure returns a PureBiConsumer.
func (c BiConsumer[T, U]) Pure() PureBiConsumer[T, U] {
	return func(t T, u U) error {
		_, err := c(context.Background(), t, u)
		return err
	}
}

// Pure returns a PureMustBiConsumer.
func (c MustBiConsumer[T, U]) Pure() PureMustBiConsumer[T, U] {
	return func(t T, u U) {
		_ = c(context.Background(), t, u)
		return
	}
}

// Must returns a PureMustBiConsumer.
func (c PureBiConsumer[T, U]) Must() PureMustBiConsumer[T, U] {
	return func(t T, u U) {
		err := c(t, u)
		if err != nil {
			panic(err)
		}
		return
	}
}

// Then returns a composed BiConsumer that performs, in sequence, this operation followed by the after operation.
func (c BiConsumer[T, U]) Then(after BiConsumer[T, U]) BiConsumer[T, U] {
	return func(ctx context.Context, t T, u U) (context.Context, error) {
		ctx, err := c(ctx, t, u)
		if err != nil {
			return ctx, err
		}
		return after(ctx, t, u)
	}
}

// Then returns a composed MustBiConsumer that performs, in sequence, this operation followed by the after operation.
func (c MustBiConsumer[T, U]) Then(after MustBiConsumer[T, U]) MustBiConsumer[T, U] {
	return func(ctx context.Context, t T, u U) context.Context {
		return after(c(ctx, t, u), t, u)
	}
}

// Then returns a composed PureBiConsumer that performs, in sequence, this operation followed by the after operation.
func (c PureBiConsumer[T, U]) Then(after PureBiConsumer[T, U]) PureBiConsumer[T, U] {
	return func(t T, u U) error {
		err := c(t, u)
		if err != nil {
			return err
		}
		return after(t, u)
	}
}

// Then returns a composed PureMustBiConsumer that performs, in sequence, this operation followed by the after operation.
func (c PureMustBiConsumer[T, U]) Then(after PureMustBiConsumer[T, U]) PureMustBiConsumer[T, U] {
	return func(t T, u U) {
		c(t, u)
		after(t, u)
	}
}
