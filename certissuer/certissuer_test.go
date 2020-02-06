package certissuer

import (
	"fmt"
	"os"
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

func TestCmdIssue(t *testing.T) {
	_ = os.Setenv("APP_ENV", "test")

	const (
		makefileDir  = ".."
		testFilepath = "/path/to/test/file/"
	)

	os.Setenv("ISSUING_SERVICE_DIR", makefileDir)
	cmd := cmdIssue(testFilepath)

	assert.Equal(t, cmd.String(), fmt.Sprintf("/usr/bin/make -C %s issue CONF_PATH=%s", makefileDir, testFilepath))

	out, err := cmd.Output()

	assert.Nil(t, err)
	assert.Equal(t, string(out), fmt.Sprintf("/usr/bin/cert-issuer -c %s\n", testFilepath))
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
