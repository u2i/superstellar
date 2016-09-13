package superstellar_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSuperstellar(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Superstellar Suite")
}
