package pdfconv

import (
	"github.com/lastrust/issuing-service/utils/filesystem"
	"github.com/lastrust/issuing-service/utils/path"

	"github.com/lastrust/utils-go/logging"
)

type (
	PdfConverter interface {
		HtmlToPdf(issuer, processId string, groupId int32, filename string) error
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

func (c *pdfConverter) HtmlToPdf(issuerId, processId string, groupId int32, filename string) error {
	var (
		certificatePath  = path.UnsignedCertificateFilepath(issuerId, processId, groupId, filename)
		tempHtmlFilepath = path.HtmlTempFilepath(issuerId, filename)
		pdfFilepath      = path.PdfFilepath(issuerId, filename)
	)
	defer filesystem.Remove(tempHtmlFilepath)

	html, err := c.htmltopdf.ParseUnsignedCertificate(certificatePath)
	if err != nil {
		logging.Err().WithError(err).Info("pdfconf.HtmlToPdf | htmltopdf.ParseUnsignedCertificate")
		return err
	}

	err = c.htmltopdf.CreateTempHtmlTemplate(html, tempHtmlFilepath)
	if err != nil {
		logging.Err().WithError(err).Info("pdfconf.HtmlToPdf | htmltopdf.CreateTempHtmlTemplate")
		return err
	}

	err = c.htmltopdf.ExecPdfGenCommand(tempHtmlFilepath, pdfFilepath)
	if err != nil {
		logging.Err().WithError(err).Info("pdfconf.HtmlToPdf | htmltopdf.ExecPdfGenCommand")
		return err
	}

	return nil
}
