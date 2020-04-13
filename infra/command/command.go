package command

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Command struct{}

func New() *Command {
	return &Command{}
}

func (Command) IssueBlockchainCertificate(confPath string) ([]byte, error) {
	cmd := exec.Command(
		"make", "issue",
		"CONF_PATH="+confPath,
	)
	logrus.Debugf("[EXECUTE] cmd: %s\n", cmd.String())
	return cmd.Output()
}

func (Command) HtmlToPdf(htmlFilepath, pdfFilepath string) ([]byte, error) {
	cmd := exec.Command(
		"make", "htmltopdf",
		"HTML_FILEPATH="+htmlFilepath,
		"PDF_FILEPATH="+pdfFilepath,
	)
	logrus.Debugf("[EXECUTE] cmd: %s\n", cmd.String())
	return cmd.Output()
}
