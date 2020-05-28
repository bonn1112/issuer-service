package htmltopdf

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/lastrust/issuing-service/domain/pdfconv"
)

var (
	ErrDisplayHTMLNotFound = errors.New("displayHtml field not found")
	ErrDisplayHTMLStruct   = errors.New("displayHtml field must be string")
)

var LayoutFilepath = "static/layout.html"

type HtmlToPdf struct {
	Command     pdfconv.Command
	Certificate map[string]interface{}
}

func New(command pdfconv.Command) *HtmlToPdf {
	return &HtmlToPdf{
		Command:     command,
		Certificate: make(map[string]interface{}),
	}
}

type layoutData struct {
	Content template.HTML
}

func (h2p *HtmlToPdf) ParseUnsignedCertificate(certPath string) (html interface{}, err error) {
	certContent, err := ioutil.ReadFile(certPath)
	if err != nil {
		return
	}

	err = json.Unmarshal(certContent, &h2p.Certificate)
	if err != nil {
		return
	}

	html, ok := h2p.Certificate["displayHtml"]
	if !ok {
		return nil, ErrDisplayHTMLNotFound
	}

	return html, nil
}

func (h2p *HtmlToPdf) CreateTempHtmlTemplate(html interface{}, htmlFilepath string) error {
	htmlString, ok := html.(string)
	if !ok {
		return ErrDisplayHTMLStruct
	}

	// TODO: rewrite to reading this file at once
	tpl, err := template.ParseFiles(LayoutFilepath)
	if err != nil {
		return fmt.Errorf("failed parsing layout file, %v", err)
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, layoutData{template.HTML(htmlString)}); err != nil {
		return fmt.Errorf("failed executing template, %v", err)
	}

	htmlFile, err := os.OpenFile(htmlFilepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		return fmt.Errorf("create temp HtmlToPdf file error, %v", err)
	}
	_, _ = htmlFile.Write(buf.Bytes())
	htmlFile.Close()
	return nil
}

func (h2p *HtmlToPdf) ExecPdfGenCommand(htmlFilepath, pdfFilepath string) error {
	return h2p.Command.HtmlToPdf(htmlFilepath, pdfFilepath)
}
