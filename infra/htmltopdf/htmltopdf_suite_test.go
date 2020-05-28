package htmltopdf_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHtmltopdf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Htmltopdf Suite")
}
