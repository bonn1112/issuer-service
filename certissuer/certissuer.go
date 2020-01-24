package certissuer

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/lastrust/issuing-service/utils"
)

const defaultCertIssuerExecutable = "/usr/bin/cert-issuer"

// A CertIssuer for issuing the blockchain certificates
type CertIssuer interface {
	// IssueCertificate using the unsigned certificate with configuration file
	// for issuing a blockchain certificate
	IssueCertificate() error
}

type certIssuer struct {
	issuer   string
	filename string

	certIssuerExecutable string
}

func (i *certIssuer) IssueCertificate() error {
	if i.filename == "" {
		return errors.New("filename couldn't be empty")
	}

	fp := i.configsFilepath()
	if !utils.FileExists(fp) {
		return errors.New("configuration file is not exists")
	}

	cmd := exec.Command(
		"make", "issue",
		"CONF_PATH="+fp,
	)
	out, err := cmd.Output()
	fmt.Println("COMMAND:", cmd.String())
	fmt.Println("OUTPUT:", string(out))
	if err != nil {
		return err
	}

	os.Remove(fp)
	os.RemoveAll(i.unsignedCertificatesDir() + i.filename)

	return nil
}

// New a certIssuer constructor
func New(issuer, fn string) CertIssuer {
	certIssuerExecutable := os.Getenv("CERT_ISSUER_EXECUTABLE")
	if certIssuerExecutable == "" {
		certIssuerExecutable = defaultCertIssuerExecutable
	}

	return &certIssuer{issuer, fn, certIssuerExecutable}
}
