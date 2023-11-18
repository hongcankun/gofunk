package funk_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	funk "github.com/hongcankun/gofunk"
)

var _ = Describe("Consumer", func() {
	var c funk.Consumer[string]
	BeforeEach(func() {
		c = func(ctx context.Context, s string) (context.Context, error) {
			return ctx, nil
		}
	})

	Describe("Converting to MustConsumer", func() {
		It("should be converted to MustConsumer", func() {
			Expect(c.Must()).To(BeAssignableToTypeOf(funk.MustConsumer[string](nil)))
		})

		When("Original consumer will return error", func() {
			BeforeEach(func() {
				c = func(ctx context.Context, s string) (context.Context, error) {
					return ctx, errors.New("")
				}
			})

			It("should panic", func() {
				Expect(func() { c.Must()(context.Background(), "") }).Should(Panic())
			})
		})

		When("Original consumer will not return error", func() {
			BeforeEach(func() {
				c = func(ctx context.Context, s string) (context.Context, error) {
					return context.WithValue(ctx, "k", 1), nil
				}
			})

			It("should return same result except error", func() {
				ctx := c.Must()(context.Background(), "")
				Expect(ctx.Value("k")).To(Equal(1))
			})
		})
	})

	Describe("Converting to PureConsumer", func() {
		BeforeEach(func() {
			c = func(ctx context.Context, s string) (context.Context, error) {
				return ctx, errors.New("")
			}
		})

		It("should be converted to PureConsumer", func() {
			Expect(c.Pure()).To(BeAssignableToTypeOf(funk.PureConsumer[string](nil)))
		})
		It("should return same result except context", func() {
			err := c.Pure()("")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Composing with another consumer", func() {
		var after funk.Consumer[string]
		BeforeEach(func() {
			after = func(ctx context.Context, s string) (context.Context, error) {
				return ctx, nil
			}
		})

		When("Consumer consumes successfully", func() {
			BeforeEach(func() {
				c = func(ctx context.Context, s string) (context.Context, error) {
					return incCtxValue(ctx), nil
				}
			})

			When("Another consumer consumes successfully", func() {
				BeforeEach(func() {
					after = func(ctx context.Context, s string) (context.Context, error) {
						return incCtxValue(ctx), nil
					}
				})

				It("should not return error", func() {
					_, err := c.Then(after)(context.Background(), "")
					Expect(err).To(Not(HaveOccurred()))
				})
				It("should execute both consumers and propagate context", func() {
					ctx, _ := c.Then(after)(context.Background(), "")
					Expect(getCtxValue(ctx)).To(Equal(2))
				})
			})

			When("Another consumer consumes unsuccessfully", func() {
				BeforeEach(func() {
					after = func(ctx context.Context, s string) (context.Context, error) {
						return incCtxValue(ctx), errors.New("")
					}
				})

				It("should return error", func() {
					_, err := c.Then(after)(context.Background(), "")
					Expect(err).To(HaveOccurred())
				})
				It("should execute both consumers and propagate context", func() {
					ctx, _ := c.Then(after)(context.Background(), "")
					Expect(getCtxValue(ctx)).To(Equal(2))
				})
			})
		})

		When("Consumer consumes unsuccessfully", func() {
			BeforeEach(func() {
				c = func(ctx context.Context, s string) (context.Context, error) {
					return incCtxValue(ctx), errors.New("")
				}
			})

			When("Another consumer consumes successfully", func() {
				BeforeEach(func() {
					after = func(ctx context.Context, s string) (context.Context, error) {
						return incCtxValue(ctx), nil
					}
				})

				It("should return error", func() {
					_, err := c.Then(after)(context.Background(), "")
					Expect(err).To(HaveOccurred())
				})
				It("should short-circuit and only execute first consumer", func() {
					ctx, _ := c.Then(after)(context.Background(), "")
					Expect(getCtxValue(ctx)).To(Equal(1))
				})
			})

			When("Another consumer consumes unsuccessfully", func() {
				BeforeEach(func() {
					after = func(ctx context.Context, s string) (context.Context, error) {
						return incCtxValue(ctx), errors.New("")
					}
				})

				It("should return error", func() {
					_, err := c.Then(after)(context.Background(), "")
					Expect(err).To(HaveOccurred())
				})
				It("should short-circuit and only execute first consumer", func() {
					ctx, _ := c.Then(after)(context.Background(), "")
					Expect(getCtxValue(ctx)).To(Equal(1))
				})
			})
		})
	})
})

var _ = Describe("MustConsumer", func() {
	var c funk.MustConsumer[string]
	var calls int
	BeforeEach(func() {
		c = func(ctx context.Context, s string) context.Context {
			return ctx
		}
		calls = 0
	})

	Describe("Converting to PureMustConsumer", func() {
		BeforeEach(func() {
			c = func(ctx context.Context, s string) context.Context {
				calls++
				return ctx
			}
		})

		It("should be converted to PureMustConsumer", func() {
			Expect(c.Pure()).To(BeAssignableToTypeOf(funk.PureMustConsumer[string](nil)))
		})
		It("should return same result except context", func() {
			c.Pure()("")
			Expect(calls).To(Equal(1))
		})
	})

	Describe("Composing with another consumer", func() {
		var after funk.MustConsumer[string]

		When("Consumer consumes successfully", func() {
			BeforeEach(func() {
				c = func(ctx context.Context, s string) context.Context {
					return incCtxValue(ctx)
				}
			})

			When("Another consumer consumes successfully", func() {
				BeforeEach(func() {
					after = func(ctx context.Context, s string) context.Context {
						return incCtxValue(ctx)
					}
				})

				It("should execute both consumers and propagate context", func() {
					ctx := c.Then(after)(context.Background(), "")
					Expect(getCtxValue(ctx)).To(Equal(2))
				})
			})
		})
	})
})

var _ = Describe("PureConsumer", func() {
	var c funk.PureConsumer[string]
	BeforeEach(func() {
		c = func(s string) error {
			return nil
		}
	})

	Describe("Converting to PureMustConsumer", func() {
		It("should be converted to PureMustConsumer", func() {
			Expect(c.Must()).To(BeAssignableToTypeOf(funk.PureMustConsumer[string](nil)))
		})

		When("Original consumer will return error", func() {
			BeforeEach(func() {
				c = func(s string) error {
					return errors.New("")
				}
			})

			It("should panic", func() {
				Expect(func() { c.Must()("") }).Should(Panic())
			})
		})

		When("Original consumer will not return error", func() {
			var calls int
			BeforeEach(func() {
				c = func(s string) error {
					calls++
					return nil
				}
				calls = 0
			})

			It("should return same result except error", func() {
				c.Must()("")
				Expect(calls).To(Equal(1))
			})
		})
	})

	Describe("Composing with another consumer", func() {
		var after funk.PureConsumer[string]
		var calls int
		BeforeEach(func() {
			calls = 0
		})

		When("Consumer consumes successfully", func() {
			BeforeEach(func() {
				c = func(s string) error {
					calls++
					return nil
				}
			})

			When("Another consumer consumes successfully", func() {
				BeforeEach(func() {
					after = func(s string) error {
						calls++
						return nil
					}
				})

				It("should not return error", func() {
					err := c.Then(after)("")
					Expect(err).To(Not(HaveOccurred()))
				})
				It("should execute both consumers and propagate context", func() {
					_ = c.Then(after)("")
					Expect(calls).To(Equal(2))
				})
			})

			When("Another consumer consumes unsuccessfully", func() {
				BeforeEach(func() {
					after = func(s string) error {
						calls++
						return errors.New("")
					}
				})

				It("should return error", func() {
					err := c.Then(after)("")
					Expect(err).To(HaveOccurred())
				})
				It("should execute both consumers and propagate context", func() {
					_ = c.Then(after)("")
					Expect(calls).To(Equal(2))
				})
			})
		})

		When("Consumer consumes unsuccessfully", func() {
			BeforeEach(func() {
				c = func(s string) error {
					calls++
					return errors.New("")
				}
			})

			When("Another consumer consumes successfully", func() {
				BeforeEach(func() {
					after = func(s string) error {
						calls++
						return nil
					}
				})

				It("should return error", func() {
					err := c.Then(after)("")
					Expect(err).To(HaveOccurred())
				})
				It("should short-circuit and only execute first consumer", func() {
					_ = c.Then(after)("")
					Expect(calls).To(Equal(1))
				})
			})

			When("Another consumer consumes unsuccessfully", func() {
				BeforeEach(func() {
					after = func(s string) error {
						calls++
						return errors.New("")
					}
				})

				It("should return error", func() {
					err := c.Then(after)("")
					Expect(err).To(HaveOccurred())
				})
				It("should short-circuit and only execute first consumer", func() {
					_ = c.Then(after)("")
					Expect(calls).To(Equal(1))
				})
			})
		})
	})
})

var _ = Describe("PureMustConsumer", func() {
	var c funk.PureMustConsumer[string]
	BeforeEach(func() {
		c = func(s string) {
			return
		}
	})

	Describe("Composing with another consumer", func() {
		var after funk.PureMustConsumer[string]
		var calls int
		BeforeEach(func() {
			calls = 0
		})

		When("Consumer consumes successfully", func() {
			BeforeEach(func() {
				c = func(s string) {
					calls++

				}
			})

			When("Another consumer consumes successfully", func() {
				BeforeEach(func() {
					after = func(s string) {
						calls++
					}
				})

				It("should execute both consumers and propagate context", func() {
					c.Then(after)("")
					Expect(calls).To(Equal(2))
				})
			})
		})
	})
})

var _ = Describe("BiConsumer", func() {
	var c funk.BiConsumer[string, string]
	BeforeEach(func() {
		c = func(ctx context.Context, _, _ string) (context.Context, error) {
			return ctx, nil
		}
	})

	Describe("Converting to MustBiConsumer", func() {
		It("should be converted to MustBiConsumer", func() {
			Expect(c.Must()).To(BeAssignableToTypeOf(funk.MustBiConsumer[string, string](nil)))
		})

		When("Original consumer will return error", func() {
			BeforeEach(func() {
				c = func(ctx context.Context, _, _ string) (context.Context, error) {
					return ctx, errors.New("")
				}
			})

			It("should panic", func() {
				Expect(func() { c.Must()(context.Background(), "", "") }).Should(Panic())
			})
		})

		When("Original consumer will not return error", func() {
			BeforeEach(func() {
				c = func(ctx context.Context, _, _ string) (context.Context, error) {
					return context.WithValue(ctx, "k", 1), nil
				}
			})

			It("should return same result except error", func() {
				ctx := c.Must()(context.Background(), "", "")
				Expect(ctx.Value("k")).To(Equal(1))
			})
		})
	})

	Describe("Converting to PureBiConsumer", func() {
		BeforeEach(func() {
			c = func(ctx context.Context, _, _ string) (context.Context, error) {
				return ctx, errors.New("")
			}
		})

		It("should be converted to PureBiConsumer", func() {
			Expect(c.Pure()).To(BeAssignableToTypeOf(funk.PureBiConsumer[string, string](nil)))
		})
		It("should return same result except context", func() {
			err := c.Pure()("", "")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Composing with another consumer", func() {
		var after funk.BiConsumer[string, string]
		BeforeEach(func() {
			after = func(ctx context.Context, _, _ string) (context.Context, error) {
				return ctx, nil
			}
		})

		When("Consumer consumes successfully", func() {
			BeforeEach(func() {
				c = func(ctx context.Context, _, _ string) (context.Context, error) {
					return incCtxValue(ctx), nil
				}
			})

			When("Another consumer consumes successfully", func() {
				BeforeEach(func() {
					after = func(ctx context.Context, _, _ string) (context.Context, error) {
						return incCtxValue(ctx), nil
					}
				})

				It("should not return error", func() {
					_, err := c.Then(after)(context.Background(), "", "")
					Expect(err).To(Not(HaveOccurred()))
				})
				It("should execute both consumers and propagate context", func() {
					ctx, _ := c.Then(after)(context.Background(), "", "")
					Expect(getCtxValue(ctx)).To(Equal(2))
				})
			})

			When("Another consumer consumes unsuccessfully", func() {
				BeforeEach(func() {
					after = func(ctx context.Context, _, _ string) (context.Context, error) {
						return incCtxValue(ctx), errors.New("")
					}
				})

				It("should return error", func() {
					_, err := c.Then(after)(context.Background(), "", "")
					Expect(err).To(HaveOccurred())
				})
				It("should execute both consumers and propagate context", func() {
					ctx, _ := c.Then(after)(context.Background(), "", "")
					Expect(getCtxValue(ctx)).To(Equal(2))
				})
			})
		})

		When("Consumer consumes unsuccessfully", func() {
			BeforeEach(func() {
				c = func(ctx context.Context, _, _ string) (context.Context, error) {
					return incCtxValue(ctx), errors.New("")
				}
			})

			When("Another consumer consumes successfully", func() {
				BeforeEach(func() {
					after = func(ctx context.Context, _, _ string) (context.Context, error) {
						return incCtxValue(ctx), nil
					}
				})

				It("should return error", func() {
					_, err := c.Then(after)(context.Background(), "", "")
					Expect(err).To(HaveOccurred())
				})
				It("should short-circuit and only execute first consumer", func() {
					ctx, _ := c.Then(after)(context.Background(), "", "")
					Expect(getCtxValue(ctx)).To(Equal(1))
				})
			})

			When("Another consumer consumes unsuccessfully", func() {
				BeforeEach(func() {
					after = func(ctx context.Context, _, _ string) (context.Context, error) {
						return incCtxValue(ctx), errors.New("")
					}
				})

				It("should return error", func() {
					_, err := c.Then(after)(context.Background(), "", "")
					Expect(err).To(HaveOccurred())
				})
				It("should short-circuit and only execute first consumer", func() {
					ctx, _ := c.Then(after)(context.Background(), "", "")
					Expect(getCtxValue(ctx)).To(Equal(1))
				})
			})
		})
	})
})

var _ = Describe("MustBiConsumer", func() {
	var c funk.MustBiConsumer[string, string]
	var calls int
	BeforeEach(func() {
		c = func(ctx context.Context, _, _ string) context.Context {
			return ctx
		}
		calls = 0
	})

	Describe("Converting to PureMustBiConsumer", func() {
		BeforeEach(func() {
			c = func(ctx context.Context, _, _ string) context.Context {
				calls++
				return ctx
			}
		})

		It("should be converted to PureMustBiConsumer", func() {
			Expect(c.Pure()).To(BeAssignableToTypeOf(funk.PureMustBiConsumer[string, string](nil)))
		})
		It("should return same result except context", func() {
			c.Pure()("", "")
			Expect(calls).To(Equal(1))
		})
	})

	Describe("Composing with another consumer", func() {
		var after funk.MustBiConsumer[string, string]

		When("Consumer consumes successfully", func() {
			BeforeEach(func() {
				c = func(ctx context.Context, _, _ string) context.Context {
					return incCtxValue(ctx)
				}
			})

			When("Another consumer consumes successfully", func() {
				BeforeEach(func() {
					after = func(ctx context.Context, _, _ string) context.Context {
						return incCtxValue(ctx)
					}
				})

				It("should execute both consumers and propagate context", func() {
					ctx := c.Then(after)(context.Background(), "", "")
					Expect(getCtxValue(ctx)).To(Equal(2))
				})
			})
		})
	})
})

var _ = Describe("PureBiConsumer", func() {
	var c funk.PureBiConsumer[string, string]
	BeforeEach(func() {
		c = func(_, _ string) error {
			return nil
		}
	})

	Describe("Converting to PureMustBiConsumer", func() {
		It("should be converted to PureMustBiConsumer", func() {
			Expect(c.Must()).To(BeAssignableToTypeOf(funk.PureMustBiConsumer[string, string](nil)))
		})

		When("Original consumer will return error", func() {
			BeforeEach(func() {
				c = func(_, _ string) error {
					return errors.New("")
				}
			})

			It("should panic", func() {
				Expect(func() { c.Must()("", "") }).Should(Panic())
			})
		})

		When("Original consumer will not return error", func() {
			var calls int
			BeforeEach(func() {
				c = func(_, _ string) error {
					calls++
					return nil
				}
				calls = 0
			})

			It("should return same result except error", func() {
				c.Must()("", "")
				Expect(calls).To(Equal(1))
			})
		})
	})

	Describe("Composing with another consumer", func() {
		var after funk.PureBiConsumer[string, string]
		var calls int
		BeforeEach(func() {
			calls = 0
		})

		When("Consumer consumes successfully", func() {
			BeforeEach(func() {
				c = func(_, _ string) error {
					calls++
					return nil
				}
			})

			When("Another consumer consumes successfully", func() {
				BeforeEach(func() {
					after = func(_, _ string) error {
						calls++
						return nil
					}
				})

				It("should not return error", func() {
					err := c.Then(after)("", "")
					Expect(err).To(Not(HaveOccurred()))
				})
				It("should execute both consumers and propagate context", func() {
					_ = c.Then(after)("", "")
					Expect(calls).To(Equal(2))
				})
			})

			When("Another consumer consumes unsuccessfully", func() {
				BeforeEach(func() {
					after = func(_, _ string) error {
						calls++
						return errors.New("")
					}
				})

				It("should return error", func() {
					err := c.Then(after)("", "")
					Expect(err).To(HaveOccurred())
				})
				It("should execute both consumers and propagate context", func() {
					_ = c.Then(after)("", "")
					Expect(calls).To(Equal(2))
				})
			})
		})

		When("Consumer consumes unsuccessfully", func() {
			BeforeEach(func() {
				c = func(_, _ string) error {
					calls++
					return errors.New("")
				}
			})

			When("Another consumer consumes successfully", func() {
				BeforeEach(func() {
					after = func(_, _ string) error {
						calls++
						return nil
					}
				})

				It("should return error", func() {
					err := c.Then(after)("", "")
					Expect(err).To(HaveOccurred())
				})
				It("should short-circuit and only execute first consumer", func() {
					_ = c.Then(after)("", "")
					Expect(calls).To(Equal(1))
				})
			})

			When("Another consumer consumes unsuccessfully", func() {
				BeforeEach(func() {
					after = func(_, _ string) error {
						calls++
						return errors.New("")
					}
				})

				It("should return error", func() {
					err := c.Then(after)("", "")
					Expect(err).To(HaveOccurred())
				})
				It("should short-circuit and only execute first consumer", func() {
					_ = c.Then(after)("", "")
					Expect(calls).To(Equal(1))
				})
			})
		})
	})
})

var _ = Describe("PureMustBiConsumer", func() {
	var c funk.PureMustBiConsumer[string, string]
	BeforeEach(func() {
		c = func(_, _ string) {
			return
		}
	})

	Describe("Composing with another consumer", func() {
		var after funk.PureMustBiConsumer[string, string]
		var calls int
		BeforeEach(func() {
			calls = 0
		})

		When("BiConsumer consumes successfully", func() {
			BeforeEach(func() {
				c = func(_, _ string) {
					calls++

				}
			})

			When("Another consumer consumes successfully", func() {
				BeforeEach(func() {
					after = func(_, _ string) {
						calls++
					}
				})

				It("should execute both consumers and propagate context", func() {
					c.Then(after)("", "")
					Expect(calls).To(Equal(2))
				})
			})
		})
	})
})
