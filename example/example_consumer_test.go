package example_test

import (
	"context"
	"fmt"

	funk "github.com/hongcankun/gofunk"
)

func ExampleConsumer() {
	c := funk.Consumer[string](func(ctx context.Context, s string) (context.Context, error) {
		fmt.Println(ctx.Value("k"), "1")
		ctx = context.WithValue(ctx, "k", 1)
		return ctx, nil
	}).Then(func(ctx context.Context, s string) (context.Context, error) {
		fmt.Println(ctx.Value("k"), "2")
		ctx = context.WithValue(ctx, "k", 2)
		return ctx, nil
	}).Pure().Then(func(s string) error {
		fmt.Println("3")
		return nil
	}).Must().Then(func(s string) {
		fmt.Println("4")
	})
	c("")

	// Output: <nil> 1
	// 1 2
	// 3
	// 4
}

func ExampleBiConsumer() {
	c := funk.BiConsumer[string, int](func(ctx context.Context, s string, i int) (context.Context, error) {
		fmt.Println(ctx.Value("k"), "1")
		ctx = context.WithValue(ctx, "k", 1)
		return ctx, nil
	}).Then(func(ctx context.Context, s string, i int) (context.Context, error) {
		fmt.Println(ctx.Value("k"), "2")
		ctx = context.WithValue(ctx, "k", 2)
		return ctx, nil
	}).Must().Then(func(ctx context.Context, s string, i int) context.Context {
		fmt.Println(ctx.Value("k"), "3")
		ctx = context.WithValue(ctx, "k", 3)
		return ctx
	}).Pure().Then(func(s string, i int) {
		fmt.Println("4")
	})
	c("", 0)

	// Output: <nil> 1
	// 1 2
	// 2 3
	// 4
}
