package htmltopdf

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"

	"github.com/lastrust/issuing-service/domain/pdfconv"
	"github.com/lastrust/utils-go/logging"
)

var (
	ErrDisplayHTMLNotFound = errors.New("displayHtml field not found")
	ErrDisplayHTMLStruct   = errors.New("displayHtml field must be string")
	ErrParseLayoutFile     = errors.New("failed parsing layout file")
)

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

func (htp *HtmlToPdf) ParseUnsignedCertificate(certPath string) (html interface{}, err error) {
	certContent, err := ioutil.ReadFile(certPath)
	if err != nil {
		return
	}

	err = json.Unmarshal(certContent, &htp.Certificate)
	if err != nil {
		return
	}

	html, ok := htp.Certificate["displayHtml"]
	if !ok {
		return nil, ErrDisplayHTMLNotFound
	}

	return html, nil
}

func (htp *HtmlToPdf) CreateTempHtmlTemplate(html interface{}, htmlFilepath string) error {
	htmlString, ok := html.(string)
	if !ok {
		return ErrDisplayHTMLStruct
	}

	// TODO: rewrite to reading this file at once
	tpl, err := template.ParseFiles("static/layout.HtmlToPdf")
	if err != nil {
		return ErrParseLayoutFile
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

func (htp *HtmlToPdf) ExecPdfGenCommand(htmlFilepath, pdfFilepath string) error {
	out, err := htp.Command.HtmlToPdf(htmlFilepath, pdfFilepath)
	cmdField := logrus.Fields{"cmd": "HtmlToPdf"}
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			logging.Err().WithFields(cmdField).WithField("stderr", string(ee.Stderr)).Debug("[EXECUTE]")
		}
		return fmt.Errorf("error command.HtmlToPdf execution, %#v", err)
	}
	logging.Out().WithFields(cmdField).WithField("stdout", string(out)).Debug("[EXECUTE]")
	return nil
}
