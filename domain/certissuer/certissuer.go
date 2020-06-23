package certissuer

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/lastrust/issuing-service/domain/cert"
	"github.com/lastrust/issuing-service/domain/pdfconv"
	"github.com/lastrust/issuing-service/utils/filesystem"
	"github.com/lastrust/issuing-service/utils/path"
	"github.com/lastrust/issuing-service/utils/str"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/semaphore"
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
	groupId        int32
	storageAdapter StorageAdapter
	command        Command
	pdfConverter   pdfconv.PdfConverter
	certRepo       cert.Repository
	semaphore      *semaphore.Weighted
	wg             sync.WaitGroup
	txLimit        int
}

// New a certIssuer constructor
func New(
	issuerId, processId string, groupId int32,
	storageAdapter StorageAdapter,
	command Command,
	pdfConverter pdfconv.PdfConverter,
	certRepo cert.Repository,
	semaphore *semaphore.Weighted,
	txLimit int,
) CertIssuer {
	return &certIssuer{
		issuerId:       issuerId,
		processId:      processId,
		groupId:        groupId,
		storageAdapter: storageAdapter,
		command:        command,
		pdfConverter:   pdfConverter,
		certRepo:       certRepo,
		semaphore:      semaphore,
		wg:             sync.WaitGroup{},
		txLimit:        txLimit,
	}
}

func (i *certIssuer) IssueCertificate(ctx context.Context) error {
	confPath := path.IssuerConfigPath(i.issuerId, i.processId, i.groupId)
	if !filesystem.FileExists(confPath) {
		return ErrNoConfig
	}
	defer filesystem.Remove(confPath)

	bcProcessDir := path.BlockcertsProcessDir(i.issuerId, i.processId, i.groupId)
	if !filesystem.FileExists(bcProcessDir) {
		_ = os.MkdirAll(bcProcessDir, 0755)
	}
	defer filesystem.RemoveAll(bcProcessDir)

	err := i.command.IssueBlockchainCertificate(confPath)
	if err != nil {
		return err
	}

	// TODO: specify the uuid, now it's hardcode :(
	if err = i.storeAllCerts(context.Background(), bcProcessDir, i.issuerId == orixUuid); err != nil {
		return fmt.Errorf("failed certIssuer.storeAllCerts, %v", err)
	}
	return nil
}

func (i *certIssuer) storeAllCerts(ctx context.Context, bcProcessDir string, withAuth bool) error {
	var (
		doneCh = make(chan bool)
		errCh  = make(chan error)

		seededRand      *rand.Rand
		passwordRecords = [][]string{
			{"cert_id", "raw_password"},
		}
	)

	if withAuth {
		seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	go func() {
		filepath.Walk(bcProcessDir, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}

			if info.IsDir() {
				return nil
			}

			i.wg.Add(1)
			i.semaphore.Acquire(ctx, 1)
			go func(file filesystem.File) {
				defer func() {
					i.wg.Done()
					i.semaphore.Release(1)
				}()

				var hashedPassword []byte
				var err error

				filename := filesystem.TrimExt(file.Info.Name())
				if err = i.storeCert(file, filename); err != nil {
					errCh <- err
					return
				}

				if withAuth {
					password := str.Random(seededRand, 16)
					passwordRecords = append(passwordRecords, []string{
						filename, password,
					})

					hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), 10)
					if err != nil {
						errCh <- err
						return
					}
				}

				c := &cert.Cert{
					Uuid:              filename,
					Password:          string(hashedPassword),
					AuthorizeRequired: true,
					IssuerId:          i.issuerId,
					IssuingProcessId:  i.processId,
				}
				// TODO: tmp solution, rewrite to transactions
				if err = i.certRepo.Create(c); err != nil {
					errCh <- err
					return
				}

			}(filesystem.File{Path: path, Info: info})
			return nil
		})
		i.wg.Wait()
		close(doneCh)
	}()

	for {
		select {
		case err := <-errCh:
			return err

		case <-doneCh:
			if withAuth {
				return i.storageAdapter.StorePasswordRecords(i.issuerId, i.processId, passwordRecords)
			}
			return nil
		}
	}
}

func (i *certIssuer) storeCert(file filesystem.File, filename string) (err error) {
	if err = i.pdfConverter.HtmlToPdf(i.issuerId, i.processId, i.groupId, filename); err != nil {
		return
	}

	pdfPath := path.PdfFilepath(i.issuerId, filename)
	if !filesystem.FileExists(pdfPath) {
		return fmt.Errorf("PDF file doesn't exist: %s", pdfPath)
	}

	err = i.storageAdapter.StorePdf(pdfPath, i.issuerId, filename)
	if err != nil {
		return
	}
	defer filesystem.Remove(pdfPath)

	return i.storageAdapter.StoreCertificate(file.Path, i.issuerId, file.Info.Name())
}
