package cert_issuer

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lastrust/issuing-service/utils"
	"github.com/sirupsen/logrus"
)

// A CertIssuer for issuing the blockchain certificates
type CertIssuer interface {
	// IssueCertificate using the unsigned certificate with configuration file
	// for issuing a blockchain certificate
	IssueCertificate() (string, error)
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

func (i *certIssuer) IssueCertificate() (string, error) {
	confPath := i.configsFilepath()
	defer os.Remove(confPath)

	if !utils.FileExists(confPath) {
		return "", errors.New("configuration file is not exists")
	}

	out, err := exec.Command("env", "CONF_PATH="+confPath, "make").Output()
	if err != nil {
		return "", err
	}
	// defer func() {
	// 	os.RemoveAll(i.unsignedCertificatesDir())
	// 	os.RemoveAll(i.blockchainCertificatesDir())
	// }()

	// err = i.storeAllCerts(i.blockchainCertificatesDir())
	// if err != nil {
	// 	return "", err
	// }

	return string(out), nil
}

func (i *certIssuer) storeAllCerts(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		logrus.Infof("Start storing certs: %s", path)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return i.storageAdapter.StoreCerts(path, i.issuer, i.filename)
		}
		logrus.Infof("Finish storing certs: %s", path)
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
