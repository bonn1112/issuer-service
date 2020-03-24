package cert_issuer

import "fmt"

var dataDir = "/storage/data/"

func (i *certIssuer) configsFilepath() string {
	return fmt.Sprintf("%s%s/configs/%s.ini", dataDir, i.issuer, i.filename)
}

func (i *certIssuer) pdfFilepath() string {
	return fmt.Sprintf("%s%s/pdf/%s.pdf", dataDir, i.issuer, i.filename)
}

func (i *certIssuer) unsignedCertificatesDir() string {
	return fmt.Sprintf("%s%s/unsigned_certificates/%s/", dataDir, i.issuer, i.filename)
}

func (i *certIssuer) blockchainCertificatesDir() string {
	return dataDir + i.issuer + "/blockchain_certificates/" + i.filename + "/"
}
