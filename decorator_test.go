package decorator_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/andriiyaremenko/decorator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Decorator", func() {
	var (
		service                      = &someService{s: "test"}
		anotherService               = &someOtherService{s: "fail"}
		anotherServiceComposedMethod func(*someOtherService, string) (string, error)
		logger                       *log
		validateDecorator            decorator.Scene
		logDecorator                 decorator.Scene
	)

	BeforeEach(func() {
		var err error

		validateDecorator, err = decorator.NewScene(
			&validate{},
			decorator.SceneDecor((*someService).SomeMethod,
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
			decorator.SceneDecor(anotherServiceComposedMethod,
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

		logger = &log{s: strings.Builder{}}
		logDecorator, err = decorator.NewScene(
			logger,
			decorator.SceneDecor((*someService).SomeMethod, (*log).LogSomeServiceSomeMethod),
			decorator.SceneDecor((*someOtherService).SomeMethod, (*log).LogSomeOtherServiceSomeMethod),
		)

		Expect(err).NotTo(HaveOccurred())
	})

	When("can validate input", func() {
		It("returns no error if valid", func() {
			result, err := decorator.
				MustDecorate((*someService).SomeMethod, validateDecorator, logDecorator).
				Call(service, "some ", false)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("some test"))
			Expect(logger.GetLog()).To(Equal("got: \"some \", false\t'resulted in \"some test\"\n"))
		})

		It("returns error if occurred", func() {
			_, err := decorator.
				MustDecorate((*someService).SomeMethod, validateDecorator, logDecorator).
				Call(service, "some ", true)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errors.New("some error")))
			Expect(logger.GetLog()).To(Equal("got: \"some \", true\t'resulted in error: some error\n"))
		})

		It("returns validation error from decorator", func() {
			_, err := decorator.
				MustDecorate((*someService).SomeMethod, validateDecorator, logDecorator).
				Call(service, "fail", true)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errors.New("validation failed")))
			Expect(logger.GetLog()).To(Equal("got: \"fail\", true\t'resulted in error: validation failed\n"))
		})

		It("can use composed (undeclared) method ", func() {
			_, err := decorator.
				MustDecorate(anotherServiceComposedMethod, validateDecorator).
				Call(anotherService, "some ")

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errors.New("validation failed")))
		})

		It("can use existing not decorated method", func() {
			res := decorator.
				MustDecorate((*someOtherService).SomeMethod, validateDecorator, logDecorator).
				Call(anotherService, "some ")
			Expect(res).To(Equal("some fail"))
			Expect(logger.GetLog()).To(Equal("got: \"some \"\t'resulted in \"some fail\"\n"))
		})
	})
	When("Decorate", func() {
		It("errors if is used with something other than func", func() {
			_, err := decorator.Decorate(struct{}{}, validateDecorator)

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

			_, err = decorator.Decorate(fn, otherDecorator)

			Expect(err).To(HaveOccurred())

			if !errors.Is(err, decorator.ErrWrongDecoratedType) {
				Fail(fmt.Sprintf("wrong error: %+v", err))
			}
		})
	})
	When("MustDecorate", func() {
		It("will panic on error", func() {
			Expect(func() { decorator.MustDecorate(struct{}{}, validateDecorator) }).To(Panic())
		})
	})
})
