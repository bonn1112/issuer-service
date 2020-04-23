package certissuer

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/lastrust/issuing-service/domain/pdfconv"
	"github.com/lastrust/issuing-service/utils/filesystem"
	"github.com/lastrust/issuing-service/utils/path"
	"github.com/lastrust/utils-go/logging"
	"github.com/sirupsen/logrus"
)

var ErrNoConfig = errors.New("configuration file is not exists")

type (
	// A CertIssuer for issuing the blockchain certificates
	CertIssuer interface {
		// IssueCertificate using the unsigned certificate with configuration file
		// for issuing a blockchain certificate
		IssueCertificate() error
	}

	StorageAdapter interface {
		StoreCerts(string, string, string) error
	}

	Command interface {
		IssueBlockchainCertificate(confPath string) ([]byte, error)
	}
)

type certIssuer struct {
	issuer         string
	processId      string
	storageAdapter StorageAdapter
	command        Command
	pdfConverter   pdfconv.PdfConverter
}

// New a certIssuer constructor
func New(
	issuer, processId string,
	storageAdapter StorageAdapter,
	command Command,
	pdfConverter pdfconv.PdfConverter,
) (CertIssuer, error) {
	return &certIssuer{
		issuer:         issuer,
		processId:      processId,
		storageAdapter: storageAdapter,
		command:        command,
		pdfConverter:   pdfConverter,
	}, nil
}

func (i *certIssuer) IssueCertificate() error {
	confPath := path.IssuerConfigPath(i.issuer, i.processId)
	defer os.Remove(confPath)

	if !filesystem.FileExists(confPath) {
		return ErrNoConfig
	}

	// FIXME: failed parsing layout file
	// if err := i.pdfConverter.HtmlToPdf(i.issuer, i.filename); err != nil {
	// 	return fmt.Errorf("failed pdfconv.PdfConverter.HtmlToPdf, %v", err)
	// }

	out, err := i.command.IssueBlockchainCertificate(confPath)
	cmdField := logrus.Fields{"cmd": "IssueBlockchainCertificate"}
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			logging.Err().WithFields(cmdField).WithField("stderr", string(ee.Stderr)).Debug("[EXECUTE]")
		}
		return fmt.Errorf("error command.IssueBlockchainCertificate execution, %#v", err)
	}
	logging.Out().WithFields(cmdField).WithField("stdout", string(out)).Debug("[EXECUTE]")

	bcCertsDir := path.BlockchainCertificatesDir(i.issuer)
	defer os.RemoveAll(bcCertsDir)

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
		return i.storageAdapter.StoreCerts(file.Path, i.issuer, file.Info.Name())
	}

	return nil
}
