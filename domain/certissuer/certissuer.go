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
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -destination=../../mocks/certissuer_Command.go -package=mocks github.com/lastrust/issuing-service/domain/certissuer Command

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
		StorePasswordRecords(issuerId, processId string, records [][]string) error
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
) CertIssuer {
	return &certIssuer{
		issuerId:       issuerId,
		processId:      processId,
		storageAdapter: storageAdapter,
		command:        command,
		pdfConverter:   pdfConverter,
		certRepo:       certRepo,
	}
}

func (i *certIssuer) IssueCertificate(ctx context.Context) error {
	confPath := path.IssuerConfigPath(i.issuerId, i.processId)
	if !filesystem.FileExists(confPath) {
		return ErrNoConfig
	}
	defer os.Remove(confPath)

	bcProcessDir := path.BlockcertsProcessDir(i.issuerId, i.processId)
	if !filesystem.FileExists(bcProcessDir) {
		_ = os.MkdirAll(bcProcessDir, 0755)
	}
	defer os.RemoveAll(bcProcessDir)

	err := i.command.IssueBlockchainCertificate(confPath)
	if err != nil {
		return err
	}

	files, err := filesystem.GetFiles(bcProcessDir)
	if err != nil {
		return err
	}

	// TODO: specify the uuid, now it's hardcode :(
	if i.issuerId == orixUuid {
		err = i.storeAllCertsWithPasswords(ctx, files)
	} else {
		err = i.storeAllCerts(ctx, files)
	}
	if err != nil {
		return fmt.Errorf("failed certIssuer.storeAllCerts, %v", err)
	}

	return nil
}

func (i *certIssuer) storeAllCerts(ctx context.Context, files []filesystem.File) error {
	certs := make([]*cert.Cert, 0)

	for _, file := range files {
		filename := filesystem.TrimExt(file.Info.Name())
		if err := i.storeCert(file, filename); err != nil {
			return err
		}

		certs = append(certs, &cert.Cert{
			Uuid:             filename,
			Password:         "",
			IssuerId:         i.issuerId,
			IssuingProcessId: i.processId,
		})
	}

	return i.certRepo.BulkCreate(ctx, certs)
}

func (i *certIssuer) storeAllCertsWithPasswords(ctx context.Context, files []filesystem.File) (err error) {
	var certs []*cert.Cert
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	passwordRecords := [][]string{
		{"cert_id", "raw_password"},
	}

	// TODO: rewrite to running as goroutine
	for _, file := range files {
		filename := filesystem.TrimExt(file.Info.Name())
		if err = i.storeCert(file, filename); err != nil {
			return err
		}

		password := str.Random(seededRand, 16)

		passwordRecords = append(passwordRecords, []string{
			filename, password,
		})

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			return err
		}

		certs = append(certs, &cert.Cert{
			Uuid:              filename,
			Password:          string(hashedPassword),
			AuthorizeRequired: true,
			IssuerId:          i.issuerId,
			IssuingProcessId:  i.processId,
		})
	}

	if err = i.storageAdapter.StorePasswordRecords(i.issuerId, i.processId, passwordRecords); err != nil {
		return err
	}

	return i.certRepo.BulkCreate(ctx, certs)
}

func (i *certIssuer) storeCert(file filesystem.File, filename string) (err error) {
	if err = i.pdfConverter.HtmlToPdf(i.issuerId, i.processId, filename); err != nil {
		return fmt.Errorf("failed pdfconv.PdfConverter.HtmlToPdf, %v", err)
	}

	pdfPath := path.PdfFilepath(i.issuerId, filename)
	if !filesystem.FileExists(pdfPath) {
		return fmt.Errorf("PDF file doesn't exist: %s", pdfPath)
	}

	err = i.storageAdapter.StorePdf(pdfPath, i.issuerId, filename)
	if err != nil {
		return err
	}
	defer os.Remove(pdfPath)

	return i.storageAdapter.StoreCertificate(file.Path, i.issuerId, file.Info.Name())
}
