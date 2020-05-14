package certissuer

import (
	"errors"
	"fmt"
	"os"

	"github.com/lastrust/issuing-service/domain/pdfconv"
	"github.com/lastrust/issuing-service/utils/filesystem"
	"github.com/lastrust/issuing-service/utils/path"
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
		StoreCertificate(string, string, string) error
		StorePdf(string, string, string) error
	}

	Command interface {
		IssueBlockchainCertificate(confPath string) error
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

	bcProcessDir := path.BlockcertsProcessDir(i.issuer, i.processId)
	if !filesystem.FileExists(bcProcessDir) {
		_ = os.MkdirAll(bcProcessDir, 0755)
	}
	defer os.RemoveAll(bcProcessDir)

	err := i.command.IssueBlockchainCertificate(confPath)
	if err != nil {
		return err
	}

	err = i.storeAllCerts(bcProcessDir)
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

	// TODO: rewrite to running as goroutine
	for _, file := range files {
		filenameWithoutExt := filesystem.FileNameWithoutExt(file.Info.Name())
		if err := i.pdfConverter.HtmlToPdf(i.issuer, i.processId, filenameWithoutExt); err != nil {
			return fmt.Errorf("failed pdfconv.PdfConverter.HtmlToPdf, %v", err)
		}

		pdfPath := path.PdfFilepath(i.issuer, filenameWithoutExt)
		if !filesystem.FileExists(pdfPath) {
			return fmt.Errorf("PDF file doesn't exist: %s", pdfPath)
		}

		err = i.storageAdapter.StorePdf(pdfPath, i.issuer, filenameWithoutExt)
		if err != nil {
			return err
		}
		defer os.Remove(pdfPath)

		err = i.storageAdapter.StoreCertificate(file.Path, i.issuer, file.Info.Name())
		if err != nil {
			return err
		}
	}

	return nil
}
