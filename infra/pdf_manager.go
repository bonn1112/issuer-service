package infra

import (
	"fmt"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/lastrust/issuing-service/utils/path"
)

// PdfManager responsible for PDF converting from HTML & saving to storage
type PdfManager interface {
	ToPdf(html string) error
	Save() error
	Link() string
}

// NewPdfManager constructor
func NewPdfManager(issuer, filename string) (PdfManager, error) {
	pdfgen, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	return &pdfManager{pdfgen, issuer, filename}, nil
}

type pdfManager struct {
	pdfgen   *wkhtmltopdf.PDFGenerator
	issuer   string
	filename string
}

func (m pdfManager) ToPdf(html string) error {
	reader := wkhtmltopdf.NewPageReader(
		strings.NewReader(html),
	)
	m.pdfgen.AddPage(reader)

	return m.pdfgen.Create()
}

func (m pdfManager) Save() error {
	return m.pdfgen.WriteFile(path.PdfFilepath(m.issuer, m.filename))
}

func (m pdfManager) Link() string {
	return fmt.Sprintf("/storage/issuer/%s/html/%s", m.issuer, m.filename)
}
