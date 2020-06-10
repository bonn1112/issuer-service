package certissuer_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCertissuer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Certissuer Suite")
}
