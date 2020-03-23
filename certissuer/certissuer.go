package certissuer

import (
	"errors"
	"github.com/lastrust/issuing-service/utils"
	"os"
	"os/exec"
	"path/filepath"
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

	confPath := i.configsFilepath()
	defer os.Remove(confPath)

	if !utils.FileExists(confPath) {
		return errors.New("configuration file is not exists")
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

// New a certIssuer constructor
func New(issuer, fn string) CertIssuer {
	return &certIssuer{issuer, fn}
}

func (i *certIssuer) storeAllCerts(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err = i.storeGCS(path)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
