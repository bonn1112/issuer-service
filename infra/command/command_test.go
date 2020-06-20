// Must be run from root
package command_test

import (
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
	err = cmd.HtmlToPdf(htmlFilepath, pdfFilepath)
	assert.Nil(t, err)
	assert.True(t, filesystem.FileExists(pdfFilepath))

	filesystem.Remove(pdfFilepath)
}
