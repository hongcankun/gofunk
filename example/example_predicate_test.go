package example_test

import (
	"context"
	"fmt"

	funk "github.com/hongcankun/gofunk"
)

func ExamplePredicate() {
	p := funk.Predicate[int](func(ctx context.Context, _ int) (context.Context, bool, error) {
		fmt.Println(ctx.Value("k"), 1)
		ctx = context.WithValue(ctx, "k", 1)
		return ctx, true, nil
	}).Must().And(func(ctx context.Context, _ int) (context.Context, bool) {
		fmt.Println(ctx.Value("k"), 2)
		ctx = context.WithValue(ctx, "k", 2)
		return ctx, false
	}).And(func(ctx context.Context, _ int) (context.Context, bool) {
		fmt.Println(ctx.Value("k"), 3)
		ctx = context.WithValue(ctx, "k", 3)
		return ctx, false
	}).Or(func(ctx context.Context, _ int) (context.Context, bool) {
		fmt.Println(ctx.Value("k"), 4)
		ctx = context.WithValue(ctx, "k", 4)
		return ctx, true
	}).Pure().Or(func(_ int) bool {
		fmt.Println(5)
		return false
	}).Not()
	b := p(0)
	fmt.Println(b)

	// Output: <nil> 1
	// 1 2
	// 2 4
	// false
}

func ExampleBiPredicate() {
	p := funk.BiPredicate[int, int](func(ctx context.Context, _ int, _ int) (context.Context, bool, error) {
		fmt.Println(ctx.Value("k"), 1)
		ctx = context.WithValue(ctx, "k", 1)
		return ctx, true, nil
	}).Must().And(func(ctx context.Context, _ int, _ int) (context.Context, bool) {
		fmt.Println(ctx.Value("k"), 2)
		ctx = context.WithValue(ctx, "k", 2)
		return ctx, false
	}).And(func(ctx context.Context, _ int, _ int) (context.Context, bool) {
		fmt.Println(ctx.Value("k"), 3)
		ctx = context.WithValue(ctx, "k", 3)
		return ctx, false
	}).Or(func(ctx context.Context, _ int, _ int) (context.Context, bool) {
		fmt.Println(ctx.Value("k"), 4)
		ctx = context.WithValue(ctx, "k", 4)
		return ctx, true
	}).Pure().Or(func(_ int, _ int) bool {
		fmt.Println(5)
		return false
	}).Not()
	b := p(0, 0)
	fmt.Println(b)

	// Output: <nil> 1
	// 1 2
	// 2 4
	// false
}
