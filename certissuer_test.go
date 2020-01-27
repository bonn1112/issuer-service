package main_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/lastrust/issuing-service/certissuer"
	"gopkg.in/ini.v1"
)

var (
	testIssuer  = "test-issuer"
	testStorage string
	testDataDir string
)

func initPaths() {
	var err error

	testStorage, err = filepath.Abs("test/storage/")
	testStorage += "/"
	if err != nil {
		log.Fatal(err)
	}

	testDataDir = testStorage + "data/"
	certissuer.SetDataDir(testDataDir)
	testDataDir += testIssuer + "/"
}

func TestCertIssuer_IssueCertificate(t *testing.T) {
	initPaths()

	ci := certissuer.New(testIssuer, "stub")
	err := ci.IssueCertificate()
	defer os.Remove(testStorage + "stubs/conf.ini")
	if err != nil {
		t.Error(err)
	}
}

// TODO: not yet required
func configure(t *testing.T, fn string) {
	content, err := ioutil.ReadFile(testStorage + "stubs/conf.ini")
	if err != nil {
		t.Error(err)
	}
	file, err := os.OpenFile(testDataDir+"configs/"+fn+".ini", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	file.Write(content)

	conf := ini.Empty()

	sct := conf.Section("")
	keys := []struct {
		name, value string
	}{
		{"unsigned_certificates_dir", testDataDir + "unsigned_certificates/"},
		{"blockchain_certificates_dir", testDataDir + "blockchain_certificates/"},
		{"usb_name", testStorage + "stubs/"},
	}
	for _, key := range keys {
		_, err = sct.NewKey(key.name, key.value)
		if err != nil {
			return
		}
	}

	conf.WriteTo(file)
}
