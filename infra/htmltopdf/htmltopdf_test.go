package htmltopdf_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lastrust/issuing-service/infra/htmltopdf"
	"github.com/lastrust/issuing-service/utils/filesystem"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HtmlToPDF.ParseUnsignedCertificate", func() {
	Context("when i got path of exists certificate with displayHtml field", func() {
		h2p := htmltopdf.New(nil)
		certPath, _ := filepath.Abs("../../test/cert_with_displayHtml.json")
		html, err := h2p.ParseUnsignedCertificate(certPath)
		It("error is nil, data contains expected string", func() {
			Expect(err).To(BeNil())
			Expect(html.(string)).To(Equal("<h1>Test HTML</h1>"))
		})
	})

	Context("when i got path of exists certificate without displayHtml field", func() {
		h2p := htmltopdf.New(nil)
		certPath, _ := filepath.Abs("../../test/cert_without_displayHtml.json")
		data, err := h2p.ParseUnsignedCertificate(certPath)
		It("got an error, data contains nil", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("displayHtml field not found"))
			Expect(data).To(BeNil())
		})
	})

	Context("when i got path of not exists certificate", func() {
		h2p := htmltopdf.New(nil)
		data, err := h2p.ParseUnsignedCertificate("/undefined_directory/undefined_certificate.json")
		It("got an error, data contains nil", func() {
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("open /undefined_directory/undefined_certificate.json: no such file or directory"))
			Expect(data).To(BeNil())
		})
	})
})

var _ = Describe("HtmlToPDF.CreateTempHtmlTemplate", func() {
	Context("when i give an html string as interface", func() {
		const html = "<h1>Test HTML</h1>"
		htmltopdf.LayoutFilepath = "../../static/layout.html"
		htmlFilepath, _ := filepath.Abs("../../test/htmltopdf_CreateTempHtmlTemplate.html")
		AfterEach(func() {
			htmltopdf.LayoutFilepath = "static/layout.html"
			filesystem.Remove(htmlFilepath)
		})

		h2p := htmltopdf.New(nil)
		err := h2p.CreateTempHtmlTemplate(interface{}(html), htmlFilepath)
		fmt.Println(err)
		It("error is nil and html template exists in filesystem storage", func() {
			Expect(err).To(BeNil())
			Expect(filesystem.FileExists(htmlFilepath)).To(Equal(true))
		})
	})
})
