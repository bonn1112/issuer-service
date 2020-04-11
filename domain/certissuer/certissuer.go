package certissuer

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

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/lastrust/issuing-service/utils/filesystem"
	"github.com/lastrust/issuing-service/utils/path"
)

var (
	ErrFilenameEmpty       = errors.New("filename couldn't be empty")
	ErrNoConfig            = errors.New("configuration file is not exists")
	ErrDisplayHTMLNotFound = errors.New("displayHtml field not found")
	ErrDisplayHTMLStruct   = errors.New("displayHtml field must be string")
	ErrParseLayoutFile     = errors.New("failed parsing layout file")
)

// A CertIssuer for issuing the blockchain certificates
type CertIssuer interface {
	// IssueCertificate using the unsigned certificate with configuration file
	// for issuing a blockchain certificate
	IssueCertificate() error
}

type StorageAdapter interface {
	StoreCerts(string, string, string) error
}

type certIssuer struct {
	issuer         string
	filename       string
	storageAdapter StorageAdapter
}

// New a certIssuer constructor
func New(issuer, filename string, storageAdapter StorageAdapter) (CertIssuer, error) {
	if filename == "" {
		return nil, errors.New("filename couldn't be empty")
	}
	return &certIssuer{issuer, filename, storageAdapter}, nil
}

func (i *certIssuer) IssueCertificate() error {
	if i.filename == "" {
		return ErrFilenameEmpty
	}

	confPath := path.ConfigsFilepath(i.issuer, i.filename)
	// FIXME: this method remove only one file in the case of bulk issuing
	defer os.Remove(confPath)

	if !filesystem.FileExists(confPath) {
		return ErrNoConfig
	}

	// FIXME: failed parsing layout file
	// if err := i.createPdfFile(); err != nil {
	// 	return fmt.Errorf("failed certIssuer.createPdfFile, %v", err)
	// }

	cmd := exec.Command("env", "CONF_PATH="+confPath, "make")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed command execution (%s), %v", cmd.String(), err)
	}
	logrus.Infof("command exec: %s | output: %s", cmd.String(), string(out))

	bcCertsDir := path.BlockchainCertificatesDir(i.issuer)
	// TODO: Uncomment after update the upload functions
	// defer func() {
	// 	os.RemoveAll(path.UnsignedCertificatesDir(i.issuer))
	// 	os.RemoveAll(bcCertsDir)
	// }()

	err = i.storeAllCerts(bcCertsDir)
	if err != nil {
		return fmt.Errorf("failed certIssuer.storeAllCerts, %v", err)
	}

	return nil
}

func (i *certIssuer) storeAllCerts(dir string) error {
	files, err := filesystem.GetFiles(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		return i.storageAdapter.StoreCerts(file.Path, i.issuer, i.filename)
	}

	return nil
}

type layoutData struct {
	Content template.HTML
}

func (i *certIssuer) createPdfFile() (err error) {
	certPath := fmt.Sprintf("%s%s.json", path.UnsignedCertificatesDir(i.issuer), i.filename)

	certContent, err := ioutil.ReadFile(certPath)
	if err != nil {
		return
	}

	cert := make(map[string]interface{})
	err = json.Unmarshal(certContent, &cert)
	if err != nil {
		return
	}

	html, ok := cert["displayHtml"]
	if !ok {
		return ErrDisplayHTMLNotFound
	}
	htmlString, ok := html.(string)
	if !ok {
		return ErrDisplayHTMLStruct
	}

	// TODO: rewrite to reading this file at once
	tpl, err := template.ParseFiles("static/layout.html")
	if err != nil {
		return ErrParseLayoutFile
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, layoutData{template.HTML(htmlString)}); err != nil {
		return fmt.Errorf("failed executing template, %v", err)
	}

	pdfgen, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return fmt.Errorf("failed wkhtmltopdf.NewPDFGenerator, %v", err)
	}
	pdfgen.AddPage(wkhtmltopdf.NewPageReader(bytes.NewBuffer(buf.Bytes())))
	if err = pdfgen.Create(); err != nil {
		return fmt.Errorf("failed wkhtmltopdf.PDFGenerator.Create, %v", err)
	}
	if err = pdfgen.WriteFile(path.PdfFilepath(i.issuer, i.filename)); err != nil {
		return fmt.Errorf("failed wkhtmltopdf.PDFGenerator.WriteFile, %v", err)
	}

	cert["displayPdf"] = fmt.Sprintf("/storage/issuer/%s/html/%s", i.issuer, i.filename)

	jsonCert, err := json.Marshal(&cert)
	if err != nil {
		return
	}

	certFile, err := os.OpenFile(certPath, os.O_RDWR, 0755)
	if err != nil {
		return
	}
	defer certFile.Close()

	_, err = certFile.WriteAt(jsonCert, 0)
	return err
}
