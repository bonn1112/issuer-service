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

type (
	// A CertIssuer for issuing the blockchain certificates
	CertIssuer interface {
		// IssueCertificate using the unsigned certificate with configuration file
		// for issuing a blockchain certificate
		IssueCertificate(context.Context) error
	}

	StorageAdapter interface {
		StoreCerts(string, string, string) error
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

	// FIXME: failed parsing layout file
	// if err := i.pdfConverter.HtmlToPdf(i.issuer, i.filename); err != nil {
	// 	return fmt.Errorf("failed pdfconv.PdfConverter.HtmlToPdf, %v", err)
	// }

	err := i.command.IssueBlockchainCertificate(confPath)
	if err != nil {
		return err
	}

	bcCertsDir := path.BlockchainCertificatesDir(i.issuerId)
	// FIXME: necessary create a separate blockchain certificate directory
	// 	with process id not common
	// 	because it's may duplicate a certificates if we running multiple requests
	defer func() {
		os.RemoveAll(bcCertsDir)
		os.Mkdir(bcCertsDir, 0755)
	}()

	err = i.storeAllCerts(ctx, bcCertsDir)
	if err != nil {
		return fmt.Errorf("failed certIssuer.storeAllCerts, %v", err)
	}

	return nil
}

const OrixUuid = "eee2bdaa-6927-4162-aa62-285976286d2f"

func (i *certIssuer) storeAllCerts(ctx context.Context, dir string) error {
	files, err := filesystem.GetFiles(dir)
	if err != nil {
		return err
	}

	var certs []*cert.Cert
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	var authorizeRequired bool
	// TODO: specify the uuid, now it's hardcode :(
	if i.issuerId == OrixUuid {
		authorizeRequired = true
	}

	for _, file := range files {
		// TODO: rewrite to running as goroutine
		err = i.storageAdapter.StoreCerts(file.Path, i.issuerId, file.Info.Name())
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
