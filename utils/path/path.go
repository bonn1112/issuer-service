package path

import "fmt"

var dataDir = "/storage/data/"

func ConfigsFilepath(issuer, filename string) string {
	return fmt.Sprintf("%s%s/configs/%s.ini", dataDir, issuer, filename)
}

func HtmlTempFilepath(issuer, filename string) string {
	return fmt.Sprintf("%s%s/html_tmp/%s.html", dataDir, issuer, filename)
}

func PdfFilepath(issuer, filename string) string {
	return fmt.Sprintf("%s%s/pdf/%s.pdf", dataDir, issuer, filename)
}

func UnsignedCertificatesDir(issuer string) string {
	return fmt.Sprintf("%s%s/unsigned_certificates/", dataDir, issuer)
}

func UnsignedCertificateFilepath(issuer, filename string) string {
	return fmt.Sprintf("%s%s.json", UnsignedCertificatesDir(issuer), filename)
}

func BlockchainCertificatesDir(issuer string) string {
	return fmt.Sprintf("%s%s/blockchain_certificates/", dataDir, issuer)
}

func CertsPathInGCS(issuer, filename string) string {
	return fmt.Sprintf("%s/blockchain_certificates/%s.json", issuer, filename)
}
