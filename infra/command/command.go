package command

import (
	"os"
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

func (Command) IssueBlockchainCertificate(confPath string) error {
	return run(exec.Command(
		"make", "issue",
		"CONF_PATH="+confPath,
	))
}

func (c *Command) HtmlToPdf(htmlFilepath, pdfFilepath string) error {
	return run(exec.Command(
		"make", "htmltopdf",
		"CHROME_BIN="+c.ChromeBin,
		"HTML_FILEPATH="+htmlFilepath,
		"PDF_FILEPATH="+pdfFilepath,
	))
}

func run(cmd *exec.Cmd) error {
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	logging.Out().
		WithField("cmd", cmd.String()).
		Debug("[EXECUTE]")

	return cmd.Run()
}
