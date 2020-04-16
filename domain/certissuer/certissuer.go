package certissuer

import (
	"errors"
	"fmt"
	"os"

	"github.com/lastrust/issuing-service/domain/pdfconv"
	"github.com/lastrust/issuing-service/utils/filesystem"
	"github.com/lastrust/issuing-service/utils/path"
	"github.com/sirupsen/logrus"
)

var (
	ErrFilenameEmpty = errors.New("filename couldn't be empty")
	ErrNoConfig      = errors.New("configuration file is not exists")
)

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
	filename       string
	storageAdapter StorageAdapter
	command        Command
	pdfConverter   pdfconv.PdfConverter
}

// New a certIssuer constructor
func New(
	issuer, processId, filename string,
	storageAdapter StorageAdapter,
	command Command,
	pdfConverter pdfconv.PdfConverter,
) (CertIssuer, error) {
	if filename == "" {
		return nil, errors.New("filename couldn't be empty")
	}
	return &certIssuer{
		issuer:         issuer,
		processId:      processId,
		filename:       filename,
		storageAdapter: storageAdapter,
		command:        command,
		pdfConverter:   pdfConverter,
	}, nil
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
	// if err := i.pdfConverter.HtmlToPdf(i.issuer, i.filename); err != nil {
	// 	return fmt.Errorf("failed pdfconv.PdfConverter.HtmlToPdf, %v", err)
	// }

	out, err := i.command.IssueBlockchainCertificate(confPath)
	if err != nil {
		return fmt.Errorf("error command.IssueBlockchainCertificate execution, %#v", err)
	}
	logrus.Debugf("[EXECUTE] command.IssueBlockchainCertificate, out: %s\n", string(out))

	bcCertsDir := path.BlockchainCertificatesDir(i.issuer)
	// TODO: Uncomment after update the upload functions
	//defer func() {
	//	os.RemoveAll(path.UnsignedCertificatesDir(i.issuer))
	//	os.RemoveAll(bcCertsDir)
	//}()

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
