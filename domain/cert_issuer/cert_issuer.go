package cert_issuer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/lastrust/issuing-service/utils"
	"path/filepath"

	"github.com/lastrust/issuing-service/utils"
)

var (
	ErrFilenameEmpty       = errors.New("filename couldn't be empty")
	ErrNoConfig            = errors.New("configuration file is not exists")
	ErrNoBlockchainCert    = errors.New("blockchain certificate file is not exists")
	ErrDisplayHTMLNotFound = errors.New("displayHtml field not found")
	ErrDisplayHTMLStruct   = errors.New("displayHtml field must be string")
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

// New a certIssuer constructor
func New(issuer, fn string, pdfgen *wkhtmltopdf.PDFGenerator) CertIssuer {
	return &certIssuer{issuer, fn, pdfgen}
}

type certIssuer struct {
	issuer   string
	filename string
	pdfgen   *wkhtmltopdf.PDFGenerator
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

	confPath := i.configsFilepath()
	defer os.Remove(confPath)

	if !utils.FileExists(confPath) {
		return ErrNoConfig
	}

	if err := i.createPdfFile(); err != nil {
		return err
	}

	_, err := exec.Command("env", "CONF_PATH="+confPath, "make").Output()
	if err != nil {
		return err
	}
	defer func() {
		os.RemoveAll(i.unsignedCertificatesDir())
		os.RemoveAll(i.blockchainCertificatesDir())
	}()

	err = i.storeAllCerts(i.blockchainCertificatesDir())
	if err != nil {
		return err
	}

	return nil
}

func (i *certIssuer) storeAllCerts(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return i.storageAdapter.StoreCerts(path, i.issuer, i.filename)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (i *certIssuer) createPdfFile() (err error) {
	certPath := i.unsignedCertificatesDir()
	if !utils.FileExists(certPath) {
		return ErrNoBlockchainCert
	}

	blockchainCertContent, err := ioutil.ReadFile(certPath)
	if err != nil {
		return
	}

	cert := make(map[string]interface{})
	err = json.Unmarshal(blockchainCertContent, &cert)
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

	pdfgen := &(*i.pdfgen)
	pdfgen.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(htmlString)))

	err = pdfgen.Create()
	if err != nil {
		return
	}

	err = pdfgen.WriteFile(i.pdfFilepath())
	if err != nil {
		return
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
