// Must be run from root
package command_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lastrust/issuing-service/infra/command"
	"github.com/lastrust/issuing-service/utils/filesystem"
	"github.com/stretchr/testify/assert"
)

func TestCommand_HtmlToPdf(t *testing.T) {
	issuerDir, err := filepath.Abs("test/storage/data/test-issuer")
	assert.Nil(t, err)

	var (
		htmlFilepath = issuerDir + "/html_tmp/sample.html"
		pdfFilepath  = issuerDir + "/pdf/sample.pdf"
	)

	cmd := command.New()
	_, err = cmd.HtmlToPdf(htmlFilepath, pdfFilepath)
	assert.Nil(t, err)
	assert.True(t, filesystem.FileExists(pdfFilepath))

	os.Remove(pdfFilepath)
}