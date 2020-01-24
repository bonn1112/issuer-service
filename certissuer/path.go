package certissuer

var dataDir = "/storage/data/"

//SetDataDir setter for dataDir variable
func SetDataDir(dir string) {
	dataDir = dir
}

func (i *certIssuer) configsFilepath() string {
	return dataDir + i.issuer + "/configs/" + i.filename + ".ini"
}

func (i *certIssuer) unsignedCertificatesDir() string {
	return dataDir + i.issuer + "/unsigned_certificates/" + i.filename + "/"
}
