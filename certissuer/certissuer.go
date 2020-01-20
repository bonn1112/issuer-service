package certissuer

import (
	"errors"
	"os"
	"os/exec"

	"github.com/lastrust/issuing-service/utils"
)

// A CertIssuer for issuing the blockchain certificates
type CertIssuer interface {
	// IssueCertificate using the unsigned certificate with configuration file
	// for issuing a blockchain certificate
	IssueCertificate() error
}

type certIssuer struct {
	issuer   string
	filename string
}

func (i *certIssuer) IssueCertificate() error {
	if i.filename == "" {
		return errors.New("filename couldn't be empty")
	}

	fp := i.configsFilepath()
	if !utils.FileExists(fp) {
		return errors.New("configuration file is not exists")
	}

	_, err := exec.Command("env", "CONF_PATH="+fp, "make").Output()
	if err != nil {
		return err
	}

	os.Remove(fp)
	os.RemoveAll(i.unsignedCertificatesDir() + i.filename)

	return nil
}

// New a certIssuer constructor
func New(issuer, fn string) CertIssuer {
	return &certIssuer{issuer, fn}
}