package path

import "fmt"

var dataDir = "/storage/data/"

func IssuerConfigPath(issuer, processId string) string {
	return fmt.Sprintf("%s%s/issuer-configs/%s.ini", dataDir, issuer, processId)
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
	return fmt.Sprintf("%s/blockchain_certificates/%s", issuer, filename)
}

func CertsPathInS3(issuer, filename string) string {
	return fmt.Sprintf("%s/%s", issuer, filename)
}
