package certissuer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	defer os.Remove(confPath)

	if !filesystem.FileExists(confPath) {
		return ErrNoConfig
	}

	if err := i.createPdfFile(); err != nil {
		return fmt.Errorf("failed certIssuer.createPdfFile, %v", err)
	}

	cmd := exec.Command("env", "CONF_PATH="+confPath, "make")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed command execution (%s), %v", cmd.String(), err)
	}
	logrus.Info(string(out))

	bcCertsDir := path.BlockchainCertificatesDir(i.issuer)
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

	pdfgen, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return fmt.Errorf("failed wkhtmltopdf.NewPDFGenerator, %v", err)
	}
	pdfgen.AddPage(wkhtmltopdf.NewPageReader(
		strings.NewReader(htmlString),
	))
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
