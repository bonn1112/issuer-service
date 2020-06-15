package service

import (
	"context"
	"os"

	"github.com/lastrust/issuing-service/domain/cert"
	"github.com/lastrust/issuing-service/domain/certissuer"
	"github.com/lastrust/issuing-service/domain/issuer"
	"github.com/lastrust/issuing-service/domain/pdfconv"
	"github.com/lastrust/issuing-service/infra/command"
	"github.com/lastrust/issuing-service/infra/htmltopdf"
	"github.com/lastrust/issuing-service/protocol"
	"github.com/lastrust/issuing-service/utils/dicontainer"
	"github.com/lastrust/issuing-service/utils/env"
	"github.com/lastrust/issuing-service/utils/path"
	"github.com/lastrust/utils-go/logging"
	"golang.org/x/sync/semaphore"
)

type issuingService struct {
	env       env.Service
	semaphore *semaphore.Weighted

	issuerRepo issuer.Repository
	certRepo   cert.Repository
}

func New(env env.Service, issuerRepo issuer.Repository, certRepo cert.Repository) protocol.IssuingServiceServer {
	return &issuingService{
		env:        env,
		semaphore:  semaphore.NewWeighted(env.ParallelLimit),
		issuerRepo: issuerRepo,
		certRepo:   certRepo,
	}
}

// IssueBlockchainCertificate run the command of pkg/cert-issuer, returns an error if is not success
func (s issuingService) IssueBlockchainCertificate(
	ctx context.Context,
	req *protocol.IssueBlockchainCertificateRequest,
) (*protocol.IssueBlockchainCertificateReply, error) {
	defer os.RemoveAll(path.UnsignedCertificatesDir(req.IssuerId, req.ProcessId))

	_, err := s.issuerRepo.FirstByUuid(req.IssuerId)
	if err != nil {
		logging.Err().WithError(err).Errorf("error in db request firstByUuid in issuer repo with name %s", req.IssuerId)
		return nil, err
	}

	storageAdapter, err := dicontainer.GetStorageAdapter(s.env.CloudService, s.env.ProcessEnv)
	if err != nil {
		logging.Err().WithError(err).Error("failed to build StorageAdapter")
		return nil, err
	}

	cmd := command.New()
	pdfConverter := pdfconv.New(htmltopdf.New(cmd), s.semaphore)

	ci := certissuer.New(
		req.IssuerId, req.ProcessId,
		storageAdapter,
		cmd,
		pdfConverter,
		s.certRepo,
	)

	logging.Out().Infof("Start issuing process: %s %s", req.IssuerId, req.ProcessId)
	if err = ci.IssueCertificate(ctx); err != nil {
		logging.Err().WithError(err).Error("failed cert_issuer.IssueCertificate")
		return nil, err
	}
	logging.Out().Infof("Finish issuing process: %s %s", req.IssuerId, req.ProcessId)

	return &protocol.IssueBlockchainCertificateReply{}, nil
}
