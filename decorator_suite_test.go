package decorator_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDecorator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Decorator Suite")
}
