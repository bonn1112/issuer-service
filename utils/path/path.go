package path

import "fmt"

var dataDir = "/storage/data/"

func ConfigsFilepath(issuer, filename string) string {
	return fmt.Sprintf("%s%s/configs/%s.ini", dataDir, issuer, filename)
}

func PdfFilepath(issuer, filename string) string {
	return fmt.Sprintf("%s%s/pdf/%s.pdf", dataDir, issuer, filename)
}

func UnsignedCertificatesDir(issuer, certFilename string) string {
	return fmt.Sprintf("%s%s/unsigned_certificates/%s/", dataDir, issuer, certFilename)
}

func BlockchainCertificatesDir(issuer, certFilename string) string {
	return fmt.Sprintf("%s%s/blockchain_certificates/%s/", dataDir, issuer, certFilename)
}

func CertsPathInGCS(issuer, filename string) string {
	return fmt.Sprintf("%s/blockchain_certificates/%s.json", issuer, filename)
}
