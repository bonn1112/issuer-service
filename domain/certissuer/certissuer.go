package certissuer

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/lastrust/issuing-service/domain/cert"
	"github.com/lastrust/issuing-service/domain/pdfconv"
	"github.com/lastrust/issuing-service/utils/filesystem"
	"github.com/lastrust/issuing-service/utils/path"
	"github.com/lastrust/issuing-service/utils/str"
)

var ErrNoConfig = errors.New("configuration file is not exists")

const orixUuid = "eee2bdaa-6927-4162-aa62-285976286d2f"

type (
	// A CertIssuer for issuing the blockchain certificates
	CertIssuer interface {
		// IssueCertificate using the unsigned certificate with configuration file
		// for issuing a blockchain certificate
		IssueCertificate(context.Context) error
	}

	StorageAdapter interface {
		StoreCertificate(string, string, string) error
		StorePdf(string, string, string) error
	}

	Command interface {
		IssueBlockchainCertificate(confPath string) error
	}
)

type certIssuer struct {
	issuerId       string
	processId      string
	storageAdapter StorageAdapter
	command        Command
	pdfConverter   pdfconv.PdfConverter
	certRepo       cert.Repository
}

// New a certIssuer constructor
func New(
	issuerId, processId string,
	storageAdapter StorageAdapter,
	command Command,
	pdfConverter pdfconv.PdfConverter,
	certRepo cert.Repository,
) (CertIssuer, error) {
	return &certIssuer{
		issuerId:       issuerId,
		processId:      processId,
		storageAdapter: storageAdapter,
		command:        command,
		pdfConverter:   pdfConverter,
		certRepo:       certRepo,
	}, nil
}

func (i *certIssuer) IssueCertificate(ctx context.Context) error {
	confPath := path.IssuerConfigPath(i.issuerId, i.processId)
	defer os.Remove(confPath)

	if !filesystem.FileExists(confPath) {
		return ErrNoConfig
	}

	bcProcessDir := path.BlockcertsProcessDir(i.issuerId, i.processId)
	if !filesystem.FileExists(bcProcessDir) {
		_ = os.MkdirAll(bcProcessDir, 0755)
	}
	defer os.RemoveAll(bcProcessDir)

	err := i.command.IssueBlockchainCertificate(confPath)
	if err != nil {
		return err
	}

	err = i.storeAllCerts(ctx, bcProcessDir)
	if err != nil {
		return fmt.Errorf("failed certIssuer.storeAllCerts, %v", err)
	}

	return nil
}

func (i *certIssuer) storeAllCerts(ctx context.Context, dir string) error {
	files, err := filesystem.GetFiles(dir)
	if err != nil {
		return err
	}

	var certs []*cert.Cert
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	var authorizeRequired bool
	// TODO: specify the uuid, now it's hardcode :(
	if i.issuerId == orixUuid {
		authorizeRequired = true
	}

	// TODO: rewrite to running as goroutine
	for _, file := range files {
		filenameWithoutExt := filesystem.FileNameWithoutExt(file.Info.Name())
		if err := i.pdfConverter.HtmlToPdf(i.issuerId, i.processId, filenameWithoutExt); err != nil {
			return fmt.Errorf("failed pdfconv.PdfConverter.HtmlToPdf, %v", err)
		}

		pdfPath := path.PdfFilepath(i.issuerId, filenameWithoutExt)
		if !filesystem.FileExists(pdfPath) {
			return fmt.Errorf("PDF file doesn't exist: %s", pdfPath)
		}

		err = i.storageAdapter.StorePdf(pdfPath, i.issuerId, filenameWithoutExt)
		if err != nil {
			return err
		}
		defer os.Remove(pdfPath)

		err = i.storageAdapter.StoreCertificate(file.Path, i.issuerId, file.Info.Name())
		if err != nil {
			return err
		}

		certs = append(certs, &cert.Cert{
			Uuid:              filesystem.TrimExt(file.Info.Name()),
			Password:          str.Random(seededRand, 16),
			AuthorizeRequired: authorizeRequired,
			IssuerId:          i.issuerId,
			IssuingProcessId:  i.processId,
		})
	}

	return i.certRepo.BulkCreate(ctx, certs)
}
