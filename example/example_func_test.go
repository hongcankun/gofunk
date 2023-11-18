package example_test

import (
	"context"
	"fmt"
	"strconv"

	funk "github.com/hongcankun/gofunk"
)

func ExampleFunc() {
	f := funk.Func[int, string](func(ctx context.Context, i int) (context.Context, string, error) {
		ctx = context.WithValue(ctx, "k", i)
		return ctx, strconv.Itoa(i), nil
	}).Must()

	ctx, v1 := f(context.Background(), 1)
	fmt.Println(ctx.Value("k"), v1)

	v2 := f.Pure()(2)
	fmt.Println(v2)

	// Output: 1 1
	// 2
}

func ExampleBiFunc() {
	f := funk.BiFunc[int, int, string](func(ctx context.Context, i int, i2 int) (context.Context, string, error) {
		ctx = context.WithValue(ctx, "k", i+i2)
		return ctx, strconv.Itoa(i + i2), nil
	}).Must()

	ctx, v1 := f(context.Background(), 1, 3)
	fmt.Println(ctx.Value("k"), v1)

	v2 := f.Pure()(2, 4)
	fmt.Println(v2)

	// Output: 4 4
	// 6
}

func ExampleUnary() {
	u := funk.Unary[int](func(ctx context.Context, i int) (context.Context, int, error) {
		fmt.Println(ctx.Value("k"), i)
		ctx = context.WithValue(ctx, "k", 1)
		return ctx, i + 1, nil
	}).Must().Then(func(ctx context.Context, i int) (context.Context, int) {
		fmt.Println(ctx.Value("k"), i)
		ctx = context.WithValue(ctx, "k", 2)
		return ctx, i + 2
	}).Pure()
	v := u(10)
	fmt.Println(v)

	// Output: <nil> 10
	// 1 11
	// 13
}
