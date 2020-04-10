package certissuer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"

	"github.com/lastrust/issuing-service/utils/filesystem"
	"github.com/lastrust/issuing-service/utils/path"
)

var (
	ErrFilenameEmpty       = errors.New("filename couldn't be empty")
	ErrNoConfig            = errors.New("configuration file is not exists")
	ErrDisplayHTMLNotFound = errors.New("displayHtml field not found")
	ErrDisplayHTMLStruct   = errors.New("displayHtml field must be string")
	ErrParseLayoutFile     = errors.New("failed parsing layout file")
)

type (
	// A CertIssuer for issuing the blockchain certificates
	CertIssuer interface {
		// IssueCertificate using the unsigned certificate with configuration file
		// for issuing a blockchain certificate
		IssueCertificate() error
	}

	StorageAdapter interface {
		StoreCerts(string, string, string) error
	}

	Command interface {
		HtmlToPdf(htmlFilepath, pdfFilepath string) ([]byte, error)
	}
)

type certIssuer struct {
	issuer         string
	filename       string
	storageAdapter StorageAdapter
	command        Command
}

// New a certIssuer constructor
func New(issuer, filename string, storageAdapter StorageAdapter, command Command) (CertIssuer, error) {
	if filename == "" {
		return nil, errors.New("filename couldn't be empty")
	}
	return &certIssuer{
		issuer:         issuer,
		filename:       filename,
		storageAdapter: storageAdapter,
		command:        command,
	}, nil
}

func (i *certIssuer) IssueCertificate() error {
	if i.filename == "" {
		return ErrFilenameEmpty
	}

	confPath := path.ConfigsFilepath(i.issuer, i.filename)
	// [FIXME] this method remove only one file in the case of bulk issuing
	defer os.Remove(confPath)

	if !filesystem.FileExists(confPath) {
		return ErrNoConfig
	}

	if err := i.createPdfFile(); err != nil {
		return fmt.Errorf("failed certIssuer.createPdfFile, %v", err)
	}

	cmd := exec.Command("env", "CONF_PATH="+confPath, "make")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed command execution (%s), %v", cmd.String(), err)
	}
	logrus.Infof("command exec: %s | output: %s", cmd.String(), string(out))

	bcCertsDir := path.BlockchainCertificatesDir(i.issuer)
	// TODO: Uncomment after update the upload functions
	// defer func() {
	// 	os.RemoveAll(path.UnsignedCertificatesDir(i.issuer))
	// 	os.RemoveAll(bcCertsDir)
	// }()

	err = i.storeAllCerts(bcCertsDir)
	if err != nil {
		return fmt.Errorf("failed certIssuer.storeAllCerts, %v", err)
	}

	return nil
}

func (i *certIssuer) storeAllCerts(dir string) error {
	files, err := filesystem.GetFiles(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		return i.storageAdapter.StoreCerts(file.Path, i.issuer, i.filename)
	}

	return nil
}

type layoutData struct {
	Content template.HTML
}

func (i *certIssuer) createPdfFile() error {
	var (
		err  error
		cert = make(map[string]interface{})
		html interface{}

		certPath     = fmt.Sprintf("%s%s.json", path.UnsignedCertificatesDir(i.issuer), i.filename)
		htmlFilepath = path.HtmlTempFilepath(i.issuer, i.filename)
	)

	defer os.Remove(htmlFilepath)

	// space for parsing unsigned certificate
	err = func() error {
		certContent, err := ioutil.ReadFile(certPath)
		if err != nil {
			return err
		}

		cert = make(map[string]interface{})
		err = json.Unmarshal(certContent, &cert)
		if err != nil {
			return err
		}

		var ok bool
		html, ok = cert["displayHtml"]
		if !ok {
			return ErrDisplayHTMLNotFound
		}

		return nil
	}()
	if err != nil {
		return err
	}

	// space for creating temporary html template
	err = func() error {
		htmlString, ok := html.(string)
		if !ok {
			return ErrDisplayHTMLStruct
		}

		// TODO: rewrite to reading this file at once
		tpl, err := template.ParseFiles("static/layout.html")
		if err != nil {
			return ErrParseLayoutFile
		}

		var buf bytes.Buffer
		if err = tpl.Execute(&buf, layoutData{template.HTML(htmlString)}); err != nil {
			return fmt.Errorf("failed executing template, %v", err)
		}

		htmlFile, err := os.OpenFile(htmlFilepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
		if err != nil {
			return fmt.Errorf("create temp html file error, %v", err)
		}
		_, _ = htmlFile.Write(buf.Bytes())
		htmlFile.Close()
		return nil
	}()
	if err != nil {
		return err
	}

	// space for executing a command
	err = func() error {
		out, err := i.command.HtmlToPdf(htmlFilepath, path.PdfFilepath(i.issuer, i.filename))
		if err != nil {
			return fmt.Errorf("error htmltopdf execution, %#v", err)
		}
		logrus.Debugf("[EXECUTE] out: %s\n", string(out))
		return nil
	}()
	if err != nil {
		return err
	}

	// space for updating unsigned certificate, add displayPdf field
	err = func() error {
		cert["displayPdf"] = fmt.Sprintf("/storage/issuer/%s/html/%s", i.issuer, i.filename)

		jsonCert, err := json.Marshal(&cert)
		if err != nil {
			return err
		}

		certFile, err := os.OpenFile(certPath, os.O_RDWR, 0755)
		if err != nil {
			return err
		}
		defer certFile.Close()

		_, err = certFile.WriteAt(jsonCert, 0)
		return err
	}()

	return err
}
