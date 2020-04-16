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

func UnsignedCertificatesDir(issuer, processId string) string {
	return fmt.Sprintf("%s%s/unsigned_certificates/%s/", dataDir, issuer, processId)
}

func UnsignedCertificateFilepath(issuer, processId, filename string) string {
	return fmt.Sprintf("%s%s.json", UnsignedCertificatesDir(issuer, processId), filename)
}

func BlockchainCertificatesDir(issuer string) string {
	return fmt.Sprintf("%s%s/blockchain_certificates/", dataDir, issuer)
}

func CertsPathInGCS(issuer, filename string) string {
	return fmt.Sprintf("%s/blockchain_certificates/%s.json", issuer, filename)
}
