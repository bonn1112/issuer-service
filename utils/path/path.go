package path

import "fmt"

var dataDir = "/storage/data/"

func IssuerConfigPath(issuerId, processId string) string {
	return fmt.Sprintf("%s%s/issuer-configs/%s.ini", dataDir, issuerId, processId)
}

func HtmlTempFilepath(issuerId, certId string) string {
	return fmt.Sprintf("%s%s/html_tmp/%s.html", dataDir, issuerId, certId)
}

func PdfFilepath(issuerId, certId string) string {
	return fmt.Sprintf("%s%s/pdf/%s.pdf", dataDir, issuerId, certId)
}

func UnsignedCertificatesDir(issuerId, processId string) string {
	return fmt.Sprintf("%s%s/unsigned_certificates/%s/", dataDir, issuerId, processId)
}

func UnsignedCertificateFilepath(issuerId, processId, certId string) string {
	return fmt.Sprintf("%s%s.json", UnsignedCertificatesDir(issuerId, processId), certId)
}

func BlockchainCertificatesDir(issuerId string) string {
	return fmt.Sprintf("%s%s/blockchain_certificates/", dataDir, issuerId)
}

func CertsPathInGCS(issuerId, certId string) string {
	return fmt.Sprintf("%s/blockchain_certificates/%s", issuerId, certId)
}
