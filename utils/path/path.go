package path

import "fmt"

var dataDir = "/storage/data/"

func SetDataDir(gotDataDir string) {
	dataDir = gotDataDir
}

func IssuerConfigPath(issuerId, processId string, groupId int32) string {
	return fmt.Sprintf("%s%s/issuer-configs/%s/%d.ini",
		dataDir, issuerId, processId, groupId)
}

func HtmlTempFilepath(issuerId, certId string) string {
	return fmt.Sprintf("%s%s/html_tmp/%s.html", dataDir, issuerId, certId)
}

func PdfFilepath(issuerId, certId string) string {
	return fmt.Sprintf("%s%s/pdf/%s.pdf", dataDir, issuerId, certId)
}

func UnsignedCertificatesDir(issuerId, processId string, groupId int32) string {
	return fmt.Sprintf("%s%s/unsigned_certificates/%s/%d/",
		dataDir, issuerId, processId, groupId)
}

func UnsignedCertificateFilepath(issuerId, processId string, groupId int32, certId string) string {
	return fmt.Sprintf("%s%s.json", UnsignedCertificatesDir(issuerId, processId, groupId), certId)
}

func BlockcertsProcessDir(issuerId, processId string, groupId int32) string {
	return fmt.Sprintf("%s%s/blockchain_certificates/%s/%d/",
		dataDir, issuerId, processId, groupId)
}

func CertsPathInGCS(issuerId, certId string) string {
	return fmt.Sprintf("%s/blockchain_certificates/%s", issuerId, certId)
}
