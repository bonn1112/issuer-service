package pdfconv

import (
	"os"

	"github.com/lastrust/issuing-service/utils/path"
)

type (
	PdfConverter interface {
		HtmlToPdf(issuer, processId, filename string) error
	}

	Command interface {
		HtmlToPdf(htmlFilepath, pdfFilepath string) error
	}

	HtmlToPdf interface {
		ParseUnsignedCertificate(certPath string) (html interface{}, err error)
		CreateTempHtmlTemplate(html interface{}, htmlFilepath string) error
		ExecPdfGenCommand(htmlFilepath, pdfFilepath string) error
	}
)

type pdfConverter struct {
	htmltopdf HtmlToPdf
}

func New(htmltopdf HtmlToPdf) PdfConverter {
	return &pdfConverter{htmltopdf}
}

func (c *pdfConverter) HtmlToPdf(issuerId, processId, filename string) (err error) {
	var (
		certificatePath  = path.UnsignedCertificateFilepath(issuerId, processId, filename)
		tempHtmlFilepath = path.HtmlTempFilepath(issuerId, filename)
		pdfFilepath      = path.PdfFilepath(issuerId, filename)
	)
	defer os.Remove(tempHtmlFilepath)

	html, err := c.htmltopdf.ParseUnsignedCertificate(certificatePath)
	if err != nil {
		return
	}

	err = c.htmltopdf.CreateTempHtmlTemplate(html, tempHtmlFilepath)
	if err != nil {
		return
	}

	return c.htmltopdf.ExecPdfGenCommand(tempHtmlFilepath, pdfFilepath)
}
