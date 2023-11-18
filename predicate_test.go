package funk_test

import (
	"context"
	"errors"

	funk "github.com/hongcankun/gofunk"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Predicate", func() {
	var p funk.Predicate[string]
	BeforeEach(func() {
		p = func(ctx context.Context, s string) (context.Context, bool, error) {
			return ctx, false, nil
		}
	})

	Describe("Converting to MustPredicate", func() {
		It("should be converted to a MustPredicate", func() {
			Expect(p.Must()).To(BeAssignableToTypeOf(funk.MustPredicate[string](nil)))
		})

		When("Original predicate will return error", func() {
			BeforeEach(func() {
				p = func(ctx context.Context, s string) (context.Context, bool, error) {
					return ctx, false, errors.New("")
				}
			})

			It("should panic", func() {
				Expect(func() { p.Must()(context.Background(), "") }).Should(Panic())
			})
		})

		When("Original predicate will not return error", func() {
			BeforeEach(func() {
				p = func(ctx context.Context, s string) (context.Context, bool, error) {
					return context.WithValue(ctx, "k", 1), true, nil
				}
			})
			It("should return same result without error", func() {
				ctx, v := p.Must()(context.Background(), "")
				Expect(ctx.Value("k")).To(Equal(1))
				Expect(v).To(BeTrue())
			})
		})
	})

	Describe("Converting to PurePredicate", func() {
		BeforeEach(func() {
			p = func(ctx context.Context, s string) (context.Context, bool, error) {
				return ctx, true, errors.New("")
			}
		})

		It("should be converted to PurePredicate", func() {
			Expect(p.Pure()).To(BeAssignableToTypeOf(funk.PurePredicate[string](nil)))
		})
		It("should return same results except context", func() {
			v, err := p.Pure()("")
			Expect(v).To(BeTrue())
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("Logical operations", func() {
		var another funk.Predicate[string]

		When("Predicate returns true", func() {
			BeforeEach(func() {
				p = func(ctx context.Context, s string) (context.Context, bool, error) {
					return incCtxValue(ctx), true, nil
				}
			})

			It("should return false if composed by NOT logical relation", func() {
				_, v, _ := p.Not()(context.Background(), "")
				Expect(v).To(BeFalse())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, s string) (context.Context, bool, error) {
						return incCtxValue(ctx), true, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return true", func() {
						_, v, _ := p.And(another)(context.Background(), "")
						Expect(v).To(BeTrue())
					})

					It("should verity both predicate and propagate context", func() {
						ctx, _, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						_, v, _ := p.Or(another)(context.Background(), "")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						ctx, _, _ := p.Or(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, s string) (context.Context, bool, error) {
						return incCtxValue(ctx), false, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						_, v, _ := p.And(another)(context.Background(), "")
						Expect(v).To(BeFalse())
					})

					It("should verity both predicate and propagate context", func() {
						ctx, _, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						_, v, _ := p.Or(another)(context.Background(), "")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						ctx, _, _ := p.Or(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})

			When("Another predicate returns error", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, s string) (context.Context, bool, error) {
						return incCtxValue(ctx), true, errors.New("")
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, _, err := p.And(another)(context.Background(), "")
						Expect(err).To(HaveOccurred())
					})

					It("should verity both predicate and propagate context", func() {
						ctx, _, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						_, v, _ := p.Or(another)(context.Background(), "")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						ctx, _, _ := p.Or(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})
		})

		When("Predicate returns false", func() {
			BeforeEach(func() {
				p = func(ctx context.Context, s string) (context.Context, bool, error) {
					return incCtxValue(ctx), false, nil
				}
			})

			It("should return true if composed by NOT logical relation", func() {
				_, v, _ := p.Not()(context.Background(), "")
				Expect(v).To(BeTrue())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, s string) (context.Context, bool, error) {
						return incCtxValue(ctx), true, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						_, v, _ := p.And(another)(context.Background(), "")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verity first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						_, v, _ := p.Or(another)(context.Background(), "")
						Expect(v).To(BeTrue())
					})
					It("should verify both predicates and propagate context", func() {
						ctx, _, _ := p.Or(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, s string) (context.Context, bool, error) {
						return incCtxValue(ctx), false, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						_, v, _ := p.And(another)(context.Background(), "")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verify first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return false", func() {
						_, v, _ := p.Or(another)(context.Background(), "")
						Expect(v).To(BeFalse())
					})
					It("should verify both predicates and propagate context", func() {
						ctx, _, _ := p.Or(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})
			})

			When("Another predicate returns error", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, s string) (context.Context, bool, error) {
						return incCtxValue(ctx), true, errors.New("")
					}
				})

				When("Composed by AND logical relation", func() {
					It("should not return error", func() {
						_, _, err := p.And(another)(context.Background(), "")
						Expect(err).ShouldNot(HaveOccurred())
					})

					It("should short-circuit and only verify first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, _, err := p.Or(another)(context.Background(), "")
						Expect(err).To(HaveOccurred())
					})
					It("should verify both predicates and propagate context", func() {
						ctx, _, _ := p.Or(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})
			})
		})

		When("Predicate returns error", func() {
			BeforeEach(func() {
				p = func(ctx context.Context, s string) (context.Context, bool, error) {
					return incCtxValue(ctx), true, errors.New("")
				}
			})

			It("should return error if composed by NOT logical relation", func() {
				_, _, err := p.Not()(context.Background(), "")
				Expect(err).To(HaveOccurred())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, s string) (context.Context, bool, error) {
						return incCtxValue(ctx), true, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, _, err := p.And(another)(context.Background(), "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, _, err := p.Or(another)(context.Background(), "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, s string) (context.Context, bool, error) {
						return incCtxValue(ctx), false, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, _, err := p.And(another)(context.Background(), "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, _, err := p.Or(another)(context.Background(), "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})

			When("Another predicate returns error", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, s string) (context.Context, bool, error) {
						return incCtxValue(ctx), true, errors.New("")
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, _, err := p.And(another)(context.Background(), "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, _, err := p.Or(another)(context.Background(), "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})
		})
	})
})

var _ = Describe("MustPredicate", func() {
	var p funk.MustPredicate[string]
	BeforeEach(func() {
		p = func(ctx context.Context, s string) (context.Context, bool) {
			return ctx, false
		}
	})

	Describe("Converting to PureMustPredicate", func() {
		BeforeEach(func() {
			p = func(ctx context.Context, s string) (context.Context, bool) {
				return ctx, true
			}
		})

		It("should be converted to PureMustPredicate", func() {
			Expect(p.Pure()).To(BeAssignableToTypeOf(funk.PureMustPredicate[string](nil)))
		})
		It("should return same results except context", func() {
			v := p.Pure()("")
			Expect(v).To(BeTrue())
		})
	})

	Describe("Logical operations", func() {
		var another funk.MustPredicate[string]

		When("Predicate returns true", func() {
			BeforeEach(func() {
				p = func(ctx context.Context, s string) (context.Context, bool) {
					return incCtxValue(ctx), true
				}
			})

			It("should return false if composed by NOT logical relation", func() {
				_, v := p.Not()(context.Background(), "")
				Expect(v).To(BeFalse())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, s string) (context.Context, bool) {
						return incCtxValue(ctx), true
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return true", func() {
						_, v := p.And(another)(context.Background(), "")
						Expect(v).To(BeTrue())
					})

					It("should verity both predicate and propagate context", func() {
						ctx, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						_, v := p.Or(another)(context.Background(), "")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						ctx, _ := p.Or(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, s string) (context.Context, bool) {
						return incCtxValue(ctx), false
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						_, v := p.And(another)(context.Background(), "")
						Expect(v).To(BeFalse())
					})

					It("should verity both predicate and propagate context", func() {
						ctx, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						_, v := p.Or(another)(context.Background(), "")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						ctx, _ := p.Or(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})
		})

		When("Predicate returns false", func() {
			BeforeEach(func() {
				p = func(ctx context.Context, s string) (context.Context, bool) {
					return incCtxValue(ctx), false
				}
			})

			It("should return true if composed by NOT logical relation", func() {
				_, v := p.Not()(context.Background(), "")
				Expect(v).To(BeTrue())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, s string) (context.Context, bool) {
						return incCtxValue(ctx), true
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						_, v := p.And(another)(context.Background(), "")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verity first predicate", func() {
						ctx, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						_, v := p.Or(another)(context.Background(), "")
						Expect(v).To(BeTrue())
					})
					It("should verify both predicates and propagate context", func() {
						ctx, _ := p.Or(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, s string) (context.Context, bool) {
						return incCtxValue(ctx), false
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						_, v := p.And(another)(context.Background(), "")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verify first predicate", func() {
						ctx, _ := p.And(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return false", func() {
						_, v := p.Or(another)(context.Background(), "")
						Expect(v).To(BeFalse())
					})
					It("should verify both predicates and propagate context", func() {
						ctx, _ := p.Or(another)(context.Background(), "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})
			})
		})
	})
})

var _ = Describe("PurePredicate", func() {
	var p funk.PurePredicate[string]
	BeforeEach(func() {
		p = func(s string) (bool, error) {
			return false, nil
		}
	})

	Describe("Converting to PureMustPredicate", func() {
		It("should be converted to a PureMustPredicate", func() {
			Expect(p.Must()).To(BeAssignableToTypeOf(funk.PureMustPredicate[string](nil)))
		})

		When("Original predicate will return error", func() {
			BeforeEach(func() {
				p = func(s string) (bool, error) {
					return false, errors.New("")
				}
			})

			It("should panic", func() {
				Expect(func() { p.Must()("") }).Should(Panic())
			})
		})

		When("Original predicate will not return error", func() {
			BeforeEach(func() {
				p = func(s string) (bool, error) {
					return true, nil
				}
			})
			It("should return same result without error", func() {
				v := p.Must()("")
				Expect(v).To(BeTrue())
			})
		})
	})

	Describe("Logical operations", func() {
		var another funk.PurePredicate[string]
		var calls int
		BeforeEach(func() {
			calls = 0
		})

		When("Predicate returns true", func() {
			BeforeEach(func() {
				p = func(s string) (bool, error) {
					calls++
					return true, nil
				}
			})

			It("should return false if composed by NOT logical relation", func() {
				v, _ := p.Not()("")
				Expect(v).To(BeFalse())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(s string) (bool, error) {
						calls++
						return true, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return true", func() {
						v, _ := p.And(another)("")
						Expect(v).To(BeTrue())
					})

					It("should verity both predicate", func() {
						_, _ = p.And(another)("")
						Expect(calls).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						v, _ := p.Or(another)("")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						_, _ = p.Or(another)("")
						Expect(calls).To(Equal(1))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(s string) (bool, error) {
						calls++
						return false, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						v, _ := p.And(another)("")
						Expect(v).To(BeFalse())
					})

					It("should verity both predicate", func() {
						_, _ = p.And(another)("")
						Expect(calls).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						v, _ := p.Or(another)("")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						_, _ = p.Or(another)("")
						Expect(calls).To(Equal(1))
					})
				})
			})

			When("Another predicate returns error", func() {
				BeforeEach(func() {
					another = func(s string) (bool, error) {
						calls++
						return true, errors.New("")
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, err := p.And(another)("")
						Expect(err).To(HaveOccurred())
					})

					It("should verity both predicate", func() {
						_, _ = p.And(another)("")
						Expect(calls).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						v, _ := p.Or(another)("")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						_, _ = p.Or(another)("")
						Expect(calls).To(Equal(1))
					})
				})
			})
		})

		When("Predicate returns false", func() {
			BeforeEach(func() {
				p = func(s string) (bool, error) {
					calls++
					return false, nil
				}
			})

			It("should return true if composed by NOT logical relation", func() {
				v, _ := p.Not()("")
				Expect(v).To(BeTrue())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(s string) (bool, error) {
						calls++
						return true, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						v, _ := p.And(another)("")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verity first predicate", func() {
						_, _ = p.And(another)("")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						v, _ := p.Or(another)("")
						Expect(v).To(BeTrue())
					})
					It("should verify both predicates", func() {
						_, _ = p.Or(another)("")
						Expect(calls).To(Equal(2))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(s string) (bool, error) {
						calls++
						return false, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						v, _ := p.And(another)("")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verify first predicate", func() {
						_, _ = p.And(another)("")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return false", func() {
						v, _ := p.Or(another)("")
						Expect(v).To(BeFalse())
					})
					It("should verify both predicates", func() {
						_, _ = p.Or(another)("")
						Expect(calls).To(Equal(2))
					})
				})
			})

			When("Another predicate returns error", func() {
				BeforeEach(func() {
					another = func(s string) (bool, error) {
						calls++
						return true, errors.New("")
					}
				})

				When("Composed by AND logical relation", func() {
					It("should not return error", func() {
						_, err := p.And(another)("")
						Expect(err).ShouldNot(HaveOccurred())
					})

					It("should short-circuit and only verify first predicate", func() {
						_, _ = p.And(another)("")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, err := p.Or(another)("")
						Expect(err).To(HaveOccurred())
					})
					It("should verify both predicates", func() {
						_, _ = p.Or(another)("")
						Expect(calls).To(Equal(2))
					})
				})
			})
		})

		When("Predicate returns error", func() {
			BeforeEach(func() {
				p = func(s string) (bool, error) {
					calls++
					return true, errors.New("")
				}
			})

			It("should return error if composed by NOT logical relation", func() {
				_, err := p.Not()("")
				Expect(err).To(HaveOccurred())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(s string) (bool, error) {
						calls++
						return true, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, err := p.And(another)("")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						_, _ = p.And(another)("")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, err := p.Or(another)("")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						_, _ = p.And(another)("")
						Expect(calls).To(Equal(1))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(s string) (bool, error) {
						calls++
						return false, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, err := p.And(another)("")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						_, _ = p.And(another)("")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, err := p.Or(another)("")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						_, _ = p.And(another)("")
						Expect(calls).To(Equal(1))
					})
				})
			})

			When("Another predicate returns error", func() {
				BeforeEach(func() {
					another = func(s string) (bool, error) {
						calls++
						return true, errors.New("")
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, err := p.And(another)("")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						_, _ = p.And(another)("")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, err := p.Or(another)("")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						_, _ = p.And(another)("")
						Expect(calls).To(Equal(1))
					})
				})
			})
		})
	})
})

var _ = Describe("PureMustPredicate", func() {
	var p funk.PureMustPredicate[string]
	BeforeEach(func() {
		p = func(s string) bool {
			return false
		}
	})

	Describe("Logical operations", func() {
		var another funk.PureMustPredicate[string]
		var calls int
		BeforeEach(func() {
			calls = 0
		})

		When("Predicate returns true", func() {
			BeforeEach(func() {
				p = func(s string) bool {
					calls++
					return true
				}
			})

			It("should return false if composed by NOT logical relation", func() {
				v := p.Not()("")
				Expect(v).To(BeFalse())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(s string) bool {
						calls++
						return true
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return true", func() {
						v := p.And(another)("")
						Expect(v).To(BeTrue())
					})

					It("should verity both predicate", func() {
						_ = p.And(another)("")
						Expect(calls).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						v := p.Or(another)("")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						_ = p.Or(another)("")
						Expect(calls).To(Equal(1))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(s string) bool {
						calls++
						return false
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						v := p.And(another)("")
						Expect(v).To(BeFalse())
					})

					It("should verity both predicate", func() {
						_ = p.And(another)("")
						Expect(calls).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						v := p.Or(another)("")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						_ = p.Or(another)("")
						Expect(calls).To(Equal(1))
					})
				})
			})
		})

		When("Predicate returns false", func() {
			BeforeEach(func() {
				p = func(s string) bool {
					calls++
					return false
				}
			})

			It("should return true if composed by NOT logical relation", func() {
				v := p.Not()("")
				Expect(v).To(BeTrue())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(s string) bool {
						calls++
						return true
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						v := p.And(another)("")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verity first predicate", func() {
						_ = p.And(another)("")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						v := p.Or(another)("")
						Expect(v).To(BeTrue())
					})
					It("should verify both predicates", func() {
						_ = p.Or(another)("")
						Expect(calls).To(Equal(2))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(s string) bool {
						calls++
						return false
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						v := p.And(another)("")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verify first predicate", func() {
						_ = p.And(another)("")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return false", func() {
						v := p.Or(another)("")
						Expect(v).To(BeFalse())
					})
					It("should verify both predicates", func() {
						_ = p.Or(another)("")
						Expect(calls).To(Equal(2))
					})
				})
			})
		})
	})
})

var _ = Describe("BiPredicate", func() {
	var p funk.BiPredicate[string, string]
	BeforeEach(func() {
		p = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
			return ctx, false, nil
		}
	})

	Describe("Converting to MustBiPredicate", func() {
		It("should be converted to a MustBiPredicate", func() {
			Expect(p.Must()).To(BeAssignableToTypeOf(funk.MustBiPredicate[string, string](nil)))
		})

		When("Original predicate will return error", func() {
			BeforeEach(func() {
				p = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
					return ctx, false, errors.New("")
				}
			})

			It("should panic", func() {
				Expect(func() { p.Must()(context.Background(), "", "") }).Should(Panic())
			})
		})

		When("Original predicate will not return error", func() {
			BeforeEach(func() {
				p = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
					return context.WithValue(ctx, "k", 1), true, nil
				}
			})
			It("should return same result without error", func() {
				ctx, v := p.Must()(context.Background(), "", "")
				Expect(ctx.Value("k")).To(Equal(1))
				Expect(v).To(BeTrue())
			})
		})
	})

	Describe("Converting to PureBiPredicate", func() {
		BeforeEach(func() {
			p = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
				return ctx, true, errors.New("")
			}
		})

		It("should be converted to PureBiPredicate", func() {
			Expect(p.Pure()).To(BeAssignableToTypeOf(funk.PureBiPredicate[string, string](nil)))
		})
		It("should return same results except context", func() {
			v, err := p.Pure()("", "")
			Expect(v).To(BeTrue())
			Expect(err).Should(HaveOccurred())
		})
	})

	Describe("Logical operations", func() {
		var another funk.BiPredicate[string, string]

		When("BiPredicate returns true", func() {
			BeforeEach(func() {
				p = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
					return incCtxValue(ctx), true, nil
				}
			})

			It("should return false if composed by NOT logical relation", func() {
				_, v, _ := p.Not()(context.Background(), "", "")
				Expect(v).To(BeFalse())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
						return incCtxValue(ctx), true, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return true", func() {
						_, v, _ := p.And(another)(context.Background(), "", "")
						Expect(v).To(BeTrue())
					})

					It("should verity both predicate and propagate context", func() {
						ctx, _, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						_, v, _ := p.Or(another)(context.Background(), "", "")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						ctx, _, _ := p.Or(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
						return incCtxValue(ctx), false, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						_, v, _ := p.And(another)(context.Background(), "", "")
						Expect(v).To(BeFalse())
					})

					It("should verity both predicate and propagate context", func() {
						ctx, _, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						_, v, _ := p.Or(another)(context.Background(), "", "")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						ctx, _, _ := p.Or(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})

			When("Another predicate returns error", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
						return incCtxValue(ctx), true, errors.New("")
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, _, err := p.And(another)(context.Background(), "", "")
						Expect(err).To(HaveOccurred())
					})

					It("should verity both predicate and propagate context", func() {
						ctx, _, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						_, v, _ := p.Or(another)(context.Background(), "", "")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						ctx, _, _ := p.Or(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})
		})

		When("BiPredicate returns false", func() {
			BeforeEach(func() {
				p = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
					return incCtxValue(ctx), false, nil
				}
			})

			It("should return true if composed by NOT logical relation", func() {
				_, v, _ := p.Not()(context.Background(), "", "")
				Expect(v).To(BeTrue())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
						return incCtxValue(ctx), true, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						_, v, _ := p.And(another)(context.Background(), "", "")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verity first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						_, v, _ := p.Or(another)(context.Background(), "", "")
						Expect(v).To(BeTrue())
					})
					It("should verify both predicates and propagate context", func() {
						ctx, _, _ := p.Or(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
						return incCtxValue(ctx), false, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						_, v, _ := p.And(another)(context.Background(), "", "")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verify first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return false", func() {
						_, v, _ := p.Or(another)(context.Background(), "", "")
						Expect(v).To(BeFalse())
					})
					It("should verify both predicates and propagate context", func() {
						ctx, _, _ := p.Or(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})
			})

			When("Another predicate returns error", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
						return incCtxValue(ctx), true, errors.New("")
					}
				})

				When("Composed by AND logical relation", func() {
					It("should not return error", func() {
						_, _, err := p.And(another)(context.Background(), "", "")
						Expect(err).ShouldNot(HaveOccurred())
					})

					It("should short-circuit and only verify first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, _, err := p.Or(another)(context.Background(), "", "")
						Expect(err).To(HaveOccurred())
					})
					It("should verify both predicates and propagate context", func() {
						ctx, _, _ := p.Or(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})
			})
		})

		When("BiPredicate returns error", func() {
			BeforeEach(func() {
				p = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
					return incCtxValue(ctx), true, errors.New("")
				}
			})

			It("should return error if composed by NOT logical relation", func() {
				_, _, err := p.Not()(context.Background(), "", "")
				Expect(err).To(HaveOccurred())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
						return incCtxValue(ctx), true, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, _, err := p.And(another)(context.Background(), "", "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, _, err := p.Or(another)(context.Background(), "", "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
						return incCtxValue(ctx), false, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, _, err := p.And(another)(context.Background(), "", "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, _, err := p.Or(another)(context.Background(), "", "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})

			When("Another predicate returns error", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, _, _ string) (context.Context, bool, error) {
						return incCtxValue(ctx), true, errors.New("")
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, _, err := p.And(another)(context.Background(), "", "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, _, err := p.Or(another)(context.Background(), "", "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						ctx, _, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})
		})
	})
})

var _ = Describe("MustBiPredicate", func() {
	var p funk.MustBiPredicate[string, string]
	BeforeEach(func() {
		p = func(ctx context.Context, _, _ string) (context.Context, bool) {
			return ctx, false
		}
	})

	Describe("Converting to PureMustBiPredicate", func() {
		BeforeEach(func() {
			p = func(ctx context.Context, _, _ string) (context.Context, bool) {
				return ctx, true
			}
		})

		It("should be converted to PureMustBiPredicate", func() {
			Expect(p.Pure()).To(BeAssignableToTypeOf(funk.PureMustBiPredicate[string, string](nil)))
		})
		It("should return same results except context", func() {
			v := p.Pure()("", "")
			Expect(v).To(BeTrue())
		})
	})

	Describe("Logical operations", func() {
		var another funk.MustBiPredicate[string, string]

		When("BiPredicate returns true", func() {
			BeforeEach(func() {
				p = func(ctx context.Context, _, _ string) (context.Context, bool) {
					return incCtxValue(ctx), true
				}
			})

			It("should return false if composed by NOT logical relation", func() {
				_, v := p.Not()(context.Background(), "", "")
				Expect(v).To(BeFalse())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, _, _ string) (context.Context, bool) {
						return incCtxValue(ctx), true
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return true", func() {
						_, v := p.And(another)(context.Background(), "", "")
						Expect(v).To(BeTrue())
					})

					It("should verity both predicate and propagate context", func() {
						ctx, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						_, v := p.Or(another)(context.Background(), "", "")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						ctx, _ := p.Or(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, _, _ string) (context.Context, bool) {
						return incCtxValue(ctx), false
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						_, v := p.And(another)(context.Background(), "", "")
						Expect(v).To(BeFalse())
					})

					It("should verity both predicate and propagate context", func() {
						ctx, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						_, v := p.Or(another)(context.Background(), "", "")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						ctx, _ := p.Or(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})
			})
		})

		When("BiPredicate returns false", func() {
			BeforeEach(func() {
				p = func(ctx context.Context, _, _ string) (context.Context, bool) {
					return incCtxValue(ctx), false
				}
			})

			It("should return true if composed by NOT logical relation", func() {
				_, v := p.Not()(context.Background(), "", "")
				Expect(v).To(BeTrue())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, _, _ string) (context.Context, bool) {
						return incCtxValue(ctx), true
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						_, v := p.And(another)(context.Background(), "", "")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verity first predicate", func() {
						ctx, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						_, v := p.Or(another)(context.Background(), "", "")
						Expect(v).To(BeTrue())
					})
					It("should verify both predicates and propagate context", func() {
						ctx, _ := p.Or(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(ctx context.Context, _, _ string) (context.Context, bool) {
						return incCtxValue(ctx), false
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						_, v := p.And(another)(context.Background(), "", "")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verify first predicate", func() {
						ctx, _ := p.And(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return false", func() {
						_, v := p.Or(another)(context.Background(), "", "")
						Expect(v).To(BeFalse())
					})
					It("should verify both predicates and propagate context", func() {
						ctx, _ := p.Or(another)(context.Background(), "", "")
						Expect(getCtxValue(ctx)).To(Equal(2))
					})
				})
			})
		})
	})
})

var _ = Describe("PureBiPredicate", func() {
	var p funk.PureBiPredicate[string, string]
	BeforeEach(func() {
		p = func(_, _ string) (bool, error) {
			return false, nil
		}
	})

	Describe("Converting to PureMustBiPredicate", func() {
		It("should be converted to a PureMustBiPredicate", func() {
			Expect(p.Must()).To(BeAssignableToTypeOf(funk.PureMustBiPredicate[string, string](nil)))
		})

		When("Original predicate will return error", func() {
			BeforeEach(func() {
				p = func(_, _ string) (bool, error) {
					return false, errors.New("")
				}
			})

			It("should panic", func() {
				Expect(func() { p.Must()("", "") }).Should(Panic())
			})
		})

		When("Original predicate will not return error", func() {
			BeforeEach(func() {
				p = func(_, _ string) (bool, error) {
					return true, nil
				}
			})
			It("should return same result without error", func() {
				v := p.Must()("", "")
				Expect(v).To(BeTrue())
			})
		})
	})

	Describe("Logical operations", func() {
		var another funk.PureBiPredicate[string, string]
		var calls int
		BeforeEach(func() {
			calls = 0
		})

		When("BiPredicate returns true", func() {
			BeforeEach(func() {
				p = func(_, _ string) (bool, error) {
					calls++
					return true, nil
				}
			})

			It("should return false if composed by NOT logical relation", func() {
				v, _ := p.Not()("", "")
				Expect(v).To(BeFalse())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(_, _ string) (bool, error) {
						calls++
						return true, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return true", func() {
						v, _ := p.And(another)("", "")
						Expect(v).To(BeTrue())
					})

					It("should verity both predicate", func() {
						_, _ = p.And(another)("", "")
						Expect(calls).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						v, _ := p.Or(another)("", "")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						_, _ = p.Or(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(_, _ string) (bool, error) {
						calls++
						return false, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						v, _ := p.And(another)("", "")
						Expect(v).To(BeFalse())
					})

					It("should verity both predicate", func() {
						_, _ = p.And(another)("", "")
						Expect(calls).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						v, _ := p.Or(another)("", "")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						_, _ = p.Or(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})
			})

			When("Another predicate returns error", func() {
				BeforeEach(func() {
					another = func(_, _ string) (bool, error) {
						calls++
						return true, errors.New("")
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, err := p.And(another)("", "")
						Expect(err).To(HaveOccurred())
					})

					It("should verity both predicate", func() {
						_, _ = p.And(another)("", "")
						Expect(calls).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						v, _ := p.Or(another)("", "")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						_, _ = p.Or(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})
			})
		})

		When("BiPredicate returns false", func() {
			BeforeEach(func() {
				p = func(_, _ string) (bool, error) {
					calls++
					return false, nil
				}
			})

			It("should return true if composed by NOT logical relation", func() {
				v, _ := p.Not()("", "")
				Expect(v).To(BeTrue())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(_, _ string) (bool, error) {
						calls++
						return true, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						v, _ := p.And(another)("", "")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verity first predicate", func() {
						_, _ = p.And(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						v, _ := p.Or(another)("", "")
						Expect(v).To(BeTrue())
					})
					It("should verify both predicates", func() {
						_, _ = p.Or(another)("", "")
						Expect(calls).To(Equal(2))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(_, _ string) (bool, error) {
						calls++
						return false, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						v, _ := p.And(another)("", "")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verify first predicate", func() {
						_, _ = p.And(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return false", func() {
						v, _ := p.Or(another)("", "")
						Expect(v).To(BeFalse())
					})
					It("should verify both predicates", func() {
						_, _ = p.Or(another)("", "")
						Expect(calls).To(Equal(2))
					})
				})
			})

			When("Another predicate returns error", func() {
				BeforeEach(func() {
					another = func(_, _ string) (bool, error) {
						calls++
						return true, errors.New("")
					}
				})

				When("Composed by AND logical relation", func() {
					It("should not return error", func() {
						_, err := p.And(another)("", "")
						Expect(err).ShouldNot(HaveOccurred())
					})

					It("should short-circuit and only verify first predicate", func() {
						_, _ = p.And(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, err := p.Or(another)("", "")
						Expect(err).To(HaveOccurred())
					})
					It("should verify both predicates", func() {
						_, _ = p.Or(another)("", "")
						Expect(calls).To(Equal(2))
					})
				})
			})
		})

		When("BiPredicate returns error", func() {
			BeforeEach(func() {
				p = func(_, _ string) (bool, error) {
					calls++
					return true, errors.New("")
				}
			})

			It("should return error if composed by NOT logical relation", func() {
				_, err := p.Not()("", "")
				Expect(err).To(HaveOccurred())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(_, _ string) (bool, error) {
						calls++
						return true, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, err := p.And(another)("", "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						_, _ = p.And(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, err := p.Or(another)("", "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						_, _ = p.And(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(_, _ string) (bool, error) {
						calls++
						return false, nil
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, err := p.And(another)("", "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						_, _ = p.And(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, err := p.Or(another)("", "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						_, _ = p.And(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})
			})

			When("Another predicate returns error", func() {
				BeforeEach(func() {
					another = func(_, _ string) (bool, error) {
						calls++
						return true, errors.New("")
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return error", func() {
						_, err := p.And(another)("", "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						_, _ = p.And(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return error", func() {
						_, err := p.Or(another)("", "")
						Expect(err).To(HaveOccurred())
					})
					It("should short-circuit and only verity first predicate", func() {
						_, _ = p.And(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})
			})
		})
	})
})

var _ = Describe("PureMustBiPredicate", func() {
	var p funk.PureMustBiPredicate[string, string]
	BeforeEach(func() {
		p = func(string, string) bool {
			return false
		}
	})

	Describe("Logical operations", func() {
		var another funk.PureMustBiPredicate[string, string]
		var calls int
		BeforeEach(func() {
			calls = 0
		})

		When("BiPredicate returns true", func() {
			BeforeEach(func() {
				p = func(string, string) bool {
					calls++
					return true
				}
			})

			It("should return false if composed by NOT logical relation", func() {
				v := p.Not()("", "")
				Expect(v).To(BeFalse())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(string, string) bool {
						calls++
						return true
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return true", func() {
						v := p.And(another)("", "")
						Expect(v).To(BeTrue())
					})

					It("should verity both predicate", func() {
						_ = p.And(another)("", "")
						Expect(calls).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						v := p.Or(another)("", "")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						_ = p.Or(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(string, string) bool {
						calls++
						return false
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						v := p.And(another)("", "")
						Expect(v).To(BeFalse())
					})

					It("should verity both predicate", func() {
						_ = p.And(another)("", "")
						Expect(calls).To(Equal(2))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						v := p.Or(another)("", "")
						Expect(v).To(BeTrue())
					})
					It("should short-circuit and only verify first predicate", func() {
						_ = p.Or(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})
			})
		})

		When("BiPredicate returns false", func() {
			BeforeEach(func() {
				p = func(string, string) bool {
					calls++
					return false
				}
			})

			It("should return true if composed by NOT logical relation", func() {
				v := p.Not()("", "")
				Expect(v).To(BeTrue())
			})

			When("Another predicate returns true", func() {
				BeforeEach(func() {
					another = func(string, string) bool {
						calls++
						return true
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						v := p.And(another)("", "")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verity first predicate", func() {
						_ = p.And(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return true", func() {
						v := p.Or(another)("", "")
						Expect(v).To(BeTrue())
					})
					It("should verify both predicates", func() {
						_ = p.Or(another)("", "")
						Expect(calls).To(Equal(2))
					})
				})
			})

			When("Another predicate returns false", func() {
				BeforeEach(func() {
					another = func(string, string) bool {
						calls++
						return false
					}
				})

				When("Composed by AND logical relation", func() {
					It("should return false", func() {
						v := p.And(another)("", "")
						Expect(v).To(BeFalse())
					})

					It("should short-circuit and only verify first predicate", func() {
						_ = p.And(another)("", "")
						Expect(calls).To(Equal(1))
					})
				})

				When("Composed by OR logical relation", func() {
					It("should return false", func() {
						v := p.Or(another)("", "")
						Expect(v).To(BeFalse())
					})
					It("should verify both predicates", func() {
						_ = p.Or(another)("", "")
						Expect(calls).To(Equal(2))
					})
				})
			})
		})
	})
})
