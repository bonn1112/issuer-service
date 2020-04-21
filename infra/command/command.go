package command

import (
	"os/exec"

	"github.com/lastrust/utils-go/logging"
)

type Command struct {
	ChromeBin string
}

func New() *Command {
	return &Command{
		ChromeBin: "/usr/bin/chromium-browser",
	}
}

func (Command) IssueBlockchainCertificate(confPath string) ([]byte, error) {
	cmd := exec.Command(
		"make", "issue",
		"CONF_PATH="+confPath,
	)
	logging.Out().Debugf("[EXECUTE] cmd: %s\n", cmd.String())
	return cmd.Output()
}

func (c *Command) HtmlToPdf(htmlFilepath, pdfFilepath string) ([]byte, error) {
	cmd := exec.Command(
		"make", "htmltopdf",
		"CHROME_BIN="+c.ChromeBin,
		"HTML_FILEPATH="+htmlFilepath,
		"PDF_FILEPATH="+pdfFilepath,
	)
	logging.Out().Debugf("[EXECUTE] cmd: %s\n", cmd.String())
	return cmd.Output()
}
