package funk_test

import (
	"context"
	"errors"

	funk "github.com/hongcankun/gofunk"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Supplier", func() {
	var s funk.Supplier[string]
	BeforeEach(func() {
		s = func(ctx context.Context) (context.Context, string, error) {
			return nil, "", nil
		}
	})

	Describe("Converting to PureSupplier", func() {
		It("should converted to a PureSupplier after calling Pure", func() {
			Expect(s.Pure()).To(BeAssignableToTypeOf(funk.PureSupplier[string](nil)))
		})
		It("should return same result except context", func() {
			v, err := s.Pure()()
			Expect(v).To(Equal(""))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("Converting to MustSupplier", func() {
		It("should converted to a MustSupplier after calling Must", func() {
			Expect(s.Must()).To(BeAssignableToTypeOf(funk.MustSupplier[string](nil)))
		})

		When("Original supplier will return error", func() {
			BeforeEach(func() {
				s = func(ctx context.Context) (context.Context, string, error) {
					ctx = context.WithValue(ctx, "k", 1)
					return ctx, "1", errors.New("error")
				}
			})

			It("should panic", func() {
				Expect(func() { s.Must()(context.Background()) }).Should(Panic())
			})
		})
		When("Original supplier will not return error", func() {
			BeforeEach(func() {
				s = func(ctx context.Context) (context.Context, string, error) {
					ctx = context.WithValue(ctx, "k", 1)
					return ctx, "1", nil
				}
			})

			It("should return same result except error", func() {
				ctx, v := s.Must()(context.Background())
				Expect(ctx.Value("k")).To(Equal(1))
				Expect(v).To(Equal("1"))
			})
		})
	})
})

var _ = Describe("PureSupplier", func() {
	var s funk.PureSupplier[string]
	BeforeEach(func() {
		s = func() (string, error) {
			return "", nil
		}
	})

	Describe("Converting to a PureMustSupplier", func() {
		It("should converted to a PureMustSupplier after calling Must", func() {
			Expect(s.Must()).To(BeAssignableToTypeOf(funk.PureMustSupplier[string](nil)))
		})

		When("Original supplier will return error", func() {
			BeforeEach(func() {
				s = func() (string, error) {
					return "1", errors.New("error")
				}
			})

			It("should panic", func() {
				Expect(func() { s.Must()() }).Should(Panic())
			})
		})
		When("Original supplier will not return error", func() {
			BeforeEach(func() {
				s = func() (string, error) {
					return "1", nil
				}
			})

			It("should return same result except error", func() {
				v := s.Must()()
				Expect(v).To(Equal("1"))
			})
		})
	})
})

var _ = Describe("MustSupplier", func() {
	var s funk.MustSupplier[string]
	BeforeEach(func() {
		s = func(ctx context.Context) (context.Context, string) {
			return nil, ""
		}
	})

	Describe("Converting to PureMustSupplier", func() {
		It("should converted to a PureMustSupplier after calling Pure", func() {
			Expect(s.Pure()).To(BeAssignableToTypeOf(funk.PureMustSupplier[string](nil)))
		})
		It("should return same result except context", func() {
			v := s.Pure()()
			Expect(v).To(Equal(""))
		})
	})
})
