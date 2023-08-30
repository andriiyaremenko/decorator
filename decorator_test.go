package decorator_test

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/andriiyaremenko/decorator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decorator", func() {
	var (
		service                      = &someService{s: "test"}
		anotherService               = &someOtherService{s: "fail"}
		anotherServiceComposedMethod func(*someOtherService, string) (string, error)
		validateDecorator            decorator.Scene[*validate]
	)

	BeforeEach(func() {
		var err error

		validateDecorator, err = decorator.NewScene(
			&validate{},
			decorator.Decorate((*someService).SomeMethod,
				func(
					d *validate, fn func(*someService, string, bool) (string, error),
				) func(*someService, string, bool) (string, error) {
					return func(s *someService, prefix string, fail bool) (string, error) {
						if err := d.Invalid(prefix); err != nil {
							return "", err
						}

						return fn(s, prefix, fail)
					}
				}),
			decorator.Decorate(anotherServiceComposedMethod,
				func(
					d *validate, fn func(*someOtherService, string) (string, error),
				) func(*someOtherService, string) (string, error) {
					return func(s *someOtherService, prefix string) (string, error) {
						if err := d.Invalid(s.S()); err != nil {
							return "", err
						}

						return s.SomeMethod(prefix), nil
					}
				}),
		)

		Expect(err).NotTo(HaveOccurred())
	})

	When("can validate input", func() {
		It("returns no error if valid", func() {
			result, err := decorator.MustGetCall(validateDecorator, (*someService).SomeMethod)(service, "some ", false)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("some test"))
		})

		It("returns error if occurred", func() {
			_, err := decorator.MustGetCall(validateDecorator, (*someService).SomeMethod)(service, "some ", true)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errors.New("some error")))
		})

		It("returns validation error from decorator", func() {
			_, err := decorator.MustGetCall(validateDecorator, (*someService).SomeMethod)(service, "fail", true)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errors.New("validation failed")))
		})

		It("can use composed (undeclared) method ", func() {
			_, err := decorator.MustGetCall(validateDecorator, anotherServiceComposedMethod)(anotherService, "some ")

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errors.New("validation failed")))
		})

		It("can use existing not decorated method", func() {
			res := decorator.MustGetCall(validateDecorator, (*someOtherService).SomeMethod)(anotherService, "some ")
			Expect(res).To(Equal("some fail"))
		})
	})
	When("GetCall", func() {
		It("errors if is used with something other than func", func() {
			_, err := decorator.GetCall(validateDecorator, struct{}{})

			Expect(err).To(HaveOccurred())

			if !errors.Is(err, decorator.ErrNotAFunc) {
				Fail(fmt.Sprintf("wrong error: %+v", err))
			}
		})

		It("errors if is used with custom broken options", func() {
			fn := func(a, b string) string { return a + b }
			otherDecorator, err := decorator.NewScene(
				&validate{},
				func(m map[string]any) error {
					t := reflect.TypeOf(fn)
					m[t.String()] = "some really stupid mistake"

					return nil
				},
			)

			Expect(err).NotTo(HaveOccurred())

			_, err = decorator.GetCall(otherDecorator, fn)

			Expect(err).To(HaveOccurred())

			if !errors.Is(err, decorator.ErrWrongDecoratedType) {
				Fail(fmt.Sprintf("wrong error: %+v", err))
			}
		})
	})
	When("MustGetCall", func() {
		It("will panic on error", func() {
			Expect(func() { decorator.MustGetCall(validateDecorator, struct{}{}) }).To(Panic())
		})
	})
})
