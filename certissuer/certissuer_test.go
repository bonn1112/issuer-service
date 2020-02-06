package certissuer

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testStorage string

func TestCertIssuer_errFilenameIsEmpty(t *testing.T) {
	err := New("test-issuer", "").IssueCertificate()
	assert.NotNil(t, err)
	assert.Equal(t, errFilenameIsEmpty.Error(), err.Error())
}

func TestCertIssuer_errConfigNotExists(t *testing.T) {
	err := New("test-issuer", "file_undefined").IssueCertificate()
	assert.NotNil(t, err)
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
