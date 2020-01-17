package certissuer

const dataDir = "/storage/data/"

func (i *certIssuer) configsFilepath() string {
	return dataDir + i.issuer + "/configs/" + i.filename + ".ini"
}

func (i *certIssuer) unsignedCertificatesDir() string {
	return dataDir + i.issuer + "/unsigned_certificates/" + i.filename + "/"
}
