package funk_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGofunk(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gofunk Suite")
}

func getCtxValue(ctx context.Context) int {
	v := ctx.Value("k")
	if v == nil {
		return 0
	} else {
		return v.(int)
	}
}
func incCtxValue(ctx context.Context) context.Context {
	return context.WithValue(ctx, "k", getCtxValue(ctx)+1)
}
