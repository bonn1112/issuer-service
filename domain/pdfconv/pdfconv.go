package pdfconv

import (
	"context"
	"os"

	"github.com/lastrust/issuing-service/utils/path"
	"github.com/lastrust/utils-go/logging"
	"golang.org/x/sync/semaphore"
)

type (
	PdfConverter interface {
		HtmlToPdf(issuer, processId, filename string)
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
	semaphore *semaphore.Weighted
}

func New(htmltopdf HtmlToPdf, semaphore *semaphore.Weighted) PdfConverter {
	return &pdfConverter{htmltopdf, semaphore}
}

func (c *pdfConverter) HtmlToPdf(issuerId, processId, filename string) {
	//we are expected that this process running into background
	c.semaphore.Acquire(context.Background(), 1)

	go func(issuerId, processId, filename string) {
		// FIXME: necessary update issuing status if got an error
		// 	now it's write only to log
		defer func() {
			c.semaphore.Release(1)
		}()

		var (
			certificatePath  = path.UnsignedCertificateFilepath(issuerId, processId, filename)
			tempHtmlFilepath = path.HtmlTempFilepath(issuerId, filename)
			pdfFilepath      = path.PdfFilepath(issuerId, filename)
		)
		defer os.Remove(tempHtmlFilepath)

		html, err := c.htmltopdf.ParseUnsignedCertificate(certificatePath)
		if err != nil {
			logging.Err().WithError(err).Info("pdfconf.HtmlToPdf | htmltopdf.ParseUnsignedCertificate")
			return
		}

		err = c.htmltopdf.CreateTempHtmlTemplate(html, tempHtmlFilepath)
		if err != nil {
			logging.Err().WithError(err).Info("pdfconf.HtmlToPdf | htmltopdf.CreateTempHtmlTemplate")
			return
		}

		err = c.htmltopdf.ExecPdfGenCommand(tempHtmlFilepath, pdfFilepath)
		if err != nil {
			logging.Err().WithError(err).Info("pdfconf.HtmlToPdf | htmltopdf.ExecPdfGenCommand")
		}
	}(issuerId, processId, filename)
}
