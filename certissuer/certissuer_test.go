package certissuer

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

var testStorage string

func TestCertIssuer_errFilenameIsEmpty(t *testing.T) {
	err := New("test-issuer", "").IssueCertificate()
	assert.Equal(t, errFilenameIsEmpty.Error(), err.Error())
}

func TestCertIssuer_errConfigNotExists(t *testing.T) {
	err := New("test-issuer", "file_undefined").IssueCertificate()
	assert.Equal(t, errConfigNotExists.Error(), err.Error())
}

func initTestPaths() error {
	var err error

	testStorage, err = filepath.Abs("test/storage/")
	testStorage += "/"
	if err != nil {
		return err
	}

	dataDir = testStorage + "data/"
	return nil
}
