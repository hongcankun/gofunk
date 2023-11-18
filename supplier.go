package funk

import "context"

// Supplier represents a supplier of results.
type Supplier[T any] func(context.Context) (context.Context, T, error)

// MustSupplier represents Supplier that doesn't return error.
type MustSupplier[T any] func(context.Context) (context.Context, T)

// PureSupplier represents Supplier that doesn't need context, and even meets the requirements of https://en.wikipedia.org/wiki/Pure_function.
type PureSupplier[T any] func() (T, error)

// PureMustSupplier represents Supplier that doesn't return error and doesn't need context.
type PureMustSupplier[T any] func() T

// Must return a MustSupplier.
func (s Supplier[T]) Must() MustSupplier[T] {
	return func(ctx context.Context) (context.Context, T) {
		ctx, v, err := s(ctx)
		if err != nil {
			panic(err)
		}
		return ctx, v
	}
}

// Pure returns a PureSupplier.
func (s Supplier[T]) Pure() PureSupplier[T] {
	return func() (T, error) {
		_, v, err := s(context.Background())
		return v, err
	}
}

// Pure returns a PureMustSupplier.
func (c MustSupplier[T]) Pure() PureMustSupplier[T] {
	return func() T {
		_, v := c(context.Background())
		return v
	}
}

// Must returns a PureMustSupplier.
func (c PureSupplier[T]) Must() PureMustSupplier[T] {
	return func() T {
		v, err := c()
		if err != nil {
			panic(err)
		}
		return v
	}
}
