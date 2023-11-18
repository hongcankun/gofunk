package funk_test

import (
	"context"
	"errors"

	. "github.com/onsi/gomega"

	funk "github.com/hongcankun/gofunk"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Func", func() {
	var f funk.Func[string, string]
	BeforeEach(func() {
		f = func(ctx context.Context, s string) (context.Context, string, error) {
			return ctx, s, nil
		}
	})

	Describe("Converting to MustFunc", func() {
		It("should converted to MustFunc", func() {
			Expect(f.Must()).To(BeAssignableToTypeOf(funk.MustFunc[string, string](nil)))
		})

		When("Original function will return error", func() {
			BeforeEach(func() {
				f = func(ctx context.Context, s string) (context.Context, string, error) {
					return ctx, "", errors.New("")
				}
			})

			It("should panic", func() {
				Expect(func() { f.Must()(context.Background(), "") }).Should(Panic())
			})
		})

		When("Original function will not return error", func() {
			BeforeEach(func() {
				f = func(ctx context.Context, s string) (context.Context, string, error) {
					return context.WithValue(ctx, "k", 1), "1", nil
				}
			})

			It("should return same result except error", func() {
				ctx, v := f.Must()(context.Background(), "")
				Expect(ctx.Value("k")).To(Equal(1))
				Expect(v).To(Equal("1"))
			})
		})
	})

	Describe("Converting to PureFunc", func() {
		BeforeEach(func() {
			f = func(ctx context.Context, s string) (context.Context, string, error) {
				return context.WithValue(ctx, "k", 1), "1", errors.New("")
			}
		})

		It("should converted to PureFunc", func() {
			Expect(f.Pure()).To(BeAssignableToTypeOf(funk.PureFunc[string, string](nil)))
		})
		It("should return same result except context", func() {
			v, err := f.Pure()("")
			Expect(v).To(Equal("1"))
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("MustFunc", func() {
	var f funk.MustFunc[string, string]
	BeforeEach(func() {
		f = func(ctx context.Context, s string) (context.Context, string) {
			return ctx, s
		}
	})

	Describe("Converting to PureMustFunc", func() {
		BeforeEach(func() {
			f = func(ctx context.Context, s string) (context.Context, string) {
				return context.WithValue(ctx, "k", 1), "1"
			}
		})

		It("should converted to PureMustFunc", func() {
			Expect(f.Pure()).To(BeAssignableToTypeOf(funk.PureMustFunc[string, string](nil)))
		})
		It("should return same result except context", func() {
			v := f.Pure()("")
			Expect(v).To(Equal("1"))
		})
	})
})

var _ = Describe("PureFunc", func() {
	var f funk.PureFunc[string, string]
	BeforeEach(func() {
		f = func(s string) (string, error) {
			return s, nil
		}
	})

	Describe("Converting to PureMustFunc", func() {
		It("should converted to PureMustFunc", func() {
			Expect(f.Must()).To(BeAssignableToTypeOf(funk.PureMustFunc[string, string](nil)))
		})

		When("Original function will return error", func() {
			BeforeEach(func() {
				f = func(s string) (string, error) {
					return "", errors.New("")
				}
			})

			It("should panic", func() {
				Expect(func() { f.Must()("") }).Should(Panic())
			})
		})

		When("Original function will not return error", func() {
			BeforeEach(func() {
				f = func(s string) (string, error) {
					return "1", nil
				}
			})

			It("should return same result except error", func() {
				v := f.Must()("")
				Expect(v).To(Equal("1"))
			})
		})
	})
})

var _ = Describe("BiFunc", func() {
	var f funk.BiFunc[string, string, string]
	BeforeEach(func() {
		f = func(ctx context.Context, _, _ string) (context.Context, string, error) {
			return ctx, "", nil
		}
	})

	Describe("Converting to MustBiFunc", func() {
		It("should converted to MustBiFunc", func() {
			Expect(f.Must()).To(BeAssignableToTypeOf(funk.MustBiFunc[string, string, string](nil)))
		})

		When("Original function will return error", func() {
			BeforeEach(func() {
				f = func(ctx context.Context, _, _ string) (context.Context, string, error) {
					return ctx, "", errors.New("")
				}
			})

			It("should panic", func() {
				Expect(func() { f.Must()(context.Background(), "", "") }).Should(Panic())
			})
		})

		When("Original function will not return error", func() {
			BeforeEach(func() {
				f = func(ctx context.Context, _, _ string) (context.Context, string, error) {
					return context.WithValue(ctx, "k", 1), "1", nil
				}
			})

			It("should return same result except error", func() {
				ctx, v := f.Must()(context.Background(), "", "")
				Expect(ctx.Value("k")).To(Equal(1))
				Expect(v).To(Equal("1"))
			})
		})
	})

	Describe("Converting to PureBiFunc", func() {
		BeforeEach(func() {
			f = func(ctx context.Context, _, _ string) (context.Context, string, error) {
				return context.WithValue(ctx, "k", 1), "1", errors.New("")
			}
		})

		It("should converted to PureBiFunc", func() {
			Expect(f.Pure()).To(BeAssignableToTypeOf(funk.PureBiFunc[string, string, string](nil)))
		})
		It("should return same result except context", func() {
			v, err := f.Pure()("", "")
			Expect(v).To(Equal("1"))
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("MustBiFunc", func() {
	var f funk.MustBiFunc[string, string, string]
	BeforeEach(func() {
		f = func(ctx context.Context, _, _ string) (context.Context, string) {
			return ctx, ""
		}
	})

	Describe("Converting to PureMustBiFunc", func() {
		BeforeEach(func() {
			f = func(ctx context.Context, _, _ string) (context.Context, string) {
				return context.WithValue(ctx, "k", 1), "1"
			}
		})

		It("should converted to PureMustBiFunc", func() {
			Expect(f.Pure()).To(BeAssignableToTypeOf(funk.PureMustBiFunc[string, string, string](nil)))
		})
		It("should return same result except context", func() {
			v := f.Pure()("", "")
			Expect(v).To(Equal("1"))
		})
	})
})

var _ = Describe("PureBiFunc", func() {
	var f funk.PureBiFunc[string, string, string]
	BeforeEach(func() {
		f = func(s string, _ string) (string, error) {
			return s, nil
		}
	})

	Describe("Converting to PureMustBiFunc", func() {
		It("should converted to PureMustBiFunc", func() {
			Expect(f.Must()).To(BeAssignableToTypeOf(funk.PureMustBiFunc[string, string, string](nil)))
		})

		When("Original function will return error", func() {
			BeforeEach(func() {
				f = func(s string, _ string) (string, error) {
					return "", errors.New("")
				}
			})

			It("should panic", func() {
				Expect(func() { f.Must()("", "") }).Should(Panic())
			})
		})

		When("Original function will not return error", func() {
			BeforeEach(func() {
				f = func(s string, _ string) (string, error) {
					return "1", nil
				}
			})

			It("should return same result except error", func() {
				v := f.Must()("", "")
				Expect(v).To(Equal("1"))
			})
		})
	})
})

var _ = Describe("Unary", func() {
	var f funk.Unary[string]
	BeforeEach(func() {
		f = func(ctx context.Context, s string) (context.Context, string, error) {
			return ctx, s, nil
		}
	})

	Describe("Converting to MustUnary", func() {
		It("should converted to MustUnary", func() {
			Expect(f.Must()).To(BeAssignableToTypeOf(funk.MustUnary[string](nil)))
		})

		When("Original function will return error", func() {
			BeforeEach(func() {
				f = func(ctx context.Context, s string) (context.Context, string, error) {
					return ctx, "", errors.New("")
				}
			})

			It("should panic", func() {
				Expect(func() { f.Must()(context.Background(), "") }).Should(Panic())
			})
		})

		When("Original function will not return error", func() {
			BeforeEach(func() {
				f = func(ctx context.Context, s string) (context.Context, string, error) {
					return context.WithValue(ctx, "k", 1), "1", nil
				}
			})

			It("should return same result except error", func() {
				ctx, v := f.Must()(context.Background(), "")
				Expect(ctx.Value("k")).To(Equal(1))
				Expect(v).To(Equal("1"))
			})
		})
	})

	Describe("Converting to PureUnary", func() {
		BeforeEach(func() {
			f = func(ctx context.Context, s string) (context.Context, string, error) {
				return context.WithValue(ctx, "k", 1), "1", errors.New("")
			}
		})

		It("should converted to PureUnary", func() {
			Expect(f.Pure()).To(BeAssignableToTypeOf(funk.PureUnary[string](nil)))
		})
		It("should return same result except context", func() {
			v, err := f.Pure()("")
			Expect(v).To(Equal("1"))
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Composing with another unary", func() {
		var after funk.Unary[string]

		When("Original unary will not return error", func() {
			BeforeEach(func() {
				f = func(ctx context.Context, s string) (context.Context, string, error) {
					return incCtxValue(ctx), "1", nil
				}
			})

			When("Another unary will not return error", func() {
				BeforeEach(func() {
					after = func(ctx context.Context, s string) (context.Context, string, error) {
						return incCtxValue(ctx), "2", nil
					}
				})

				It("should not return error", func() {
					_, _, err := f.Then(after)(context.Background(), "")
					Expect(err).To(Not(HaveOccurred()))
				})
				It("should execute both unary and propagate context", func() {
					ctx, v, _ := f.Then(after)(context.Background(), "")
					Expect(getCtxValue(ctx)).To(Equal(2))
					Expect(v).To(Equal("2"))
				})
			})

			When("Another unary will return error", func() {
				BeforeEach(func() {
					after = func(ctx context.Context, s string) (context.Context, string, error) {
						return incCtxValue(ctx), "", errors.New("")
					}
				})

				It("should return error", func() {
					_, _, err := f.Then(after)(context.Background(), "")
					Expect(err).To(HaveOccurred())
				})
				It("should execute both unary and propagate context", func() {
					ctx, _, _ := f.Then(after)(context.Background(), "")
					Expect(getCtxValue(ctx)).To(Equal(2))
				})
			})
		})

		When("Original unary will return error", func() {
			BeforeEach(func() {
				f = func(ctx context.Context, s string) (context.Context, string, error) {
					return incCtxValue(ctx), "", errors.New("")
				}
			})

			When("Another unary will not return error", func() {
				BeforeEach(func() {
					after = func(ctx context.Context, s string) (context.Context, string, error) {
						return incCtxValue(ctx), "2", nil
					}
				})

				It("should return error", func() {
					_, _, err := f.Then(after)(context.Background(), "")
					Expect(err).To(HaveOccurred())
				})
				It("should short-circuit and only execute first unary", func() {
					ctx, _, _ := f.Then(after)(context.Background(), "")
					Expect(getCtxValue(ctx)).To(Equal(1))
				})
			})

			When("Another unary will return error", func() {
				BeforeEach(func() {
					after = func(ctx context.Context, s string) (context.Context, string, error) {
						return incCtxValue(ctx), "", errors.New("")
					}
				})

				It("should return error", func() {
					_, _, err := f.Then(after)(context.Background(), "")
					Expect(err).To(HaveOccurred())
				})
				It("should short-circuit and only execute first unary", func() {
					ctx, _, _ := f.Then(after)(context.Background(), "")
					Expect(getCtxValue(ctx)).To(Equal(1))
				})
			})
		})
	})
})

var _ = Describe("MustUnary", func() {
	var f funk.MustUnary[string]
	BeforeEach(func() {
		f = func(ctx context.Context, s string) (context.Context, string) {
			return ctx, s
		}
	})

	Describe("Converting to PureMustUnary", func() {
		BeforeEach(func() {
			f = func(ctx context.Context, s string) (context.Context, string) {
				return context.WithValue(ctx, "k", 1), "1"
			}
		})

		It("should converted to PureMustUnary", func() {
			Expect(f.Pure()).To(BeAssignableToTypeOf(funk.PureMustUnary[string](nil)))
		})
		It("should return same result except context", func() {
			v := f.Pure()("")
			Expect(v).To(Equal("1"))
		})
	})

	Describe("Composing with another unary", func() {
		var after funk.MustUnary[string]
		BeforeEach(func() {
			f = func(ctx context.Context, s string) (context.Context, string) {
				return incCtxValue(ctx), "1"
			}
			after = func(ctx context.Context, s string) (context.Context, string) {
				return incCtxValue(ctx), "2"
			}
		})

		It("should execute both unary and propagate context", func() {
			ctx, v := f.Then(after)(context.Background(), "")
			Expect(getCtxValue(ctx)).To(Equal(2))
			Expect(v).To(Equal("2"))
		})
	})
})

var _ = Describe("PureUnary", func() {
	var f funk.PureUnary[string]
	BeforeEach(func() {
		f = func(s string) (string, error) {
			return s, nil
		}
	})

	Describe("Converting to PureMustUnary", func() {
		It("should converted to PureMustUnary", func() {
			Expect(f.Must()).To(BeAssignableToTypeOf(funk.PureMustUnary[string](nil)))
		})

		When("Original function will return error", func() {
			BeforeEach(func() {
				f = func(s string) (string, error) {
					return "", errors.New("")
				}
			})

			It("should panic", func() {
				Expect(func() { f.Must()("") }).Should(Panic())
			})
		})

		When("Original function will not return error", func() {
			BeforeEach(func() {
				f = func(s string) (string, error) {
					return "1", nil
				}
			})

			It("should return same result except error", func() {
				v := f.Must()("")
				Expect(v).To(Equal("1"))
			})
		})
	})

	Describe("Composing with another unary", func() {
		var after funk.PureUnary[string]
		var calls int
		BeforeEach(func() {
			calls = 0
		})

		When("Original unary will not return error", func() {
			BeforeEach(func() {
				f = func(s string) (string, error) {
					calls++
					return "1", nil
				}
			})

			When("Another unary will not return error", func() {
				BeforeEach(func() {
					after = func(s string) (string, error) {
						calls++
						return "2", nil
					}
				})

				It("should not return error", func() {
					_, err := f.Then(after)("")
					Expect(err).To(Not(HaveOccurred()))
				})
				It("should execute both unary and propagate context", func() {
					v, _ := f.Then(after)("")
					Expect(calls).To(Equal(2))
					Expect(v).To(Equal("2"))
				})
			})

			When("Another unary will return error", func() {
				BeforeEach(func() {
					after = func(s string) (string, error) {
						calls++
						return "", errors.New("")
					}
				})

				It("should return error", func() {
					_, err := f.Then(after)("")
					Expect(err).To(HaveOccurred())
				})
				It("should execute both unary and propagate context", func() {
					_, _ = f.Then(after)("")
					Expect(calls).To(Equal(2))
				})
			})
		})

		When("Original unary will return error", func() {
			BeforeEach(func() {
				f = func(s string) (string, error) {
					calls++
					return "", errors.New("")
				}
			})

			When("Another unary will not return error", func() {
				BeforeEach(func() {
					after = func(s string) (string, error) {
						calls++
						return "2", nil
					}
				})

				It("should return error", func() {
					_, err := f.Then(after)("")
					Expect(err).To(HaveOccurred())
				})
				It("should short-circuit and only execute first unary", func() {
					_, _ = f.Then(after)("")
					Expect(calls).To(Equal(1))
				})
			})

			When("Another unary will return error", func() {
				BeforeEach(func() {
					after = func(s string) (string, error) {
						calls++
						return "", errors.New("")
					}
				})

				It("should return error", func() {
					_, err := f.Then(after)("")
					Expect(err).To(HaveOccurred())
				})
				It("should short-circuit and only execute first unary", func() {
					_, _ = f.Then(after)("")
					Expect(calls).To(Equal(1))
				})
			})
		})
	})
})

var _ = Describe("PureMustUnary", func() {
	var f funk.PureMustUnary[string]
	BeforeEach(func() {
		f = func(s string) string {
			return s
		}
	})

	Describe("Composing with another unary", func() {
		var after funk.PureMustUnary[string]
		var calls int
		BeforeEach(func() {
			calls = 0
			f = func(s string) string {
				calls++
				return "1"
			}
			after = func(s string) string {
				calls++
				return "2"
			}
		})

		It("should execute both unary and propagate context", func() {
			v := f.Then(after)("")
			Expect(calls).To(Equal(2))
			Expect(v).To(Equal("2"))
		})
	})
})
