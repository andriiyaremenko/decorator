package decorator_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/andriiyaremenko/decorator"
)

var _ = Describe("Scene", func() {
	When("Option application returns an error", func() {
		It("NewScene will return an ErrNotAFunc error", func() {
			_, err := decorator.NewScene(
				&validate{},
				decorator.Decorate("really stupid mistake",
					func(d *validate, fn string) string { return "really stupid mistake" },
				),
			)

			Expect(err).To(HaveOccurred())
			if !errors.Is(err, decorator.ErrNotAFunc) {
				Fail(fmt.Sprintf("wrong error: %+v", err))
			}
		})

		It("NewScene will return and ErrDuplicate error", func() {
			_, err := decorator.NewScene(
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
			)

			Expect(err).To(HaveOccurred())
			if !errors.Is(err, decorator.ErrDuplicate) {
				Fail(fmt.Sprintf("wrong error: %+v", err))
			}
		})
	})
})
