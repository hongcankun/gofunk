package example_test

import (
	"context"
	"fmt"

	funk "github.com/hongcankun/gofunk"
)

func ExampleSupplier() {
	s := funk.Supplier[string](func(ctx context.Context) (context.Context, string, error) {
		ctx = context.WithValue(ctx, "k", 1)
		return ctx, "1", nil
	}).Must()

	ctx, v := s(context.Background())
	fmt.Println(ctx.Value("k"), v)

	v = s.Pure()()
	fmt.Println(v)

	// Output: 1 1
	// 1
}
