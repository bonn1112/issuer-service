package service

import (
	"context"

	"github.com/lastrust/issuing-service/domain/certissuer"
	"github.com/lastrust/issuing-service/domain/pdfconv"
	"github.com/lastrust/issuing-service/infra/command"
	"github.com/lastrust/issuing-service/infra/htmltopdf"
	"github.com/lastrust/issuing-service/protocol"
	"github.com/lastrust/issuing-service/utils/dicontainer"
	"github.com/sirupsen/logrus"
)

type issuingService struct {
	cloudService string
	processEnv   string
}

func New(cloudService, processEnv string) protocol.IssuingServiceServer {
	return &issuingService{cloudService, processEnv}
}

// IssueBlockchainCertificate run the command of pkg/cert-issuer, returns an error if is not success
func (s issuingService) IssueBlockchainCertificate(
	ctx context.Context,
	req *protocol.IssueBlockchainCertificateRequest,
) (*protocol.IssueBlockchainCertificateReply, error) {
	storageAdapter, err := dicontainer.GetStorageAdapter(s.cloudService, s.processEnv)
	if err != nil {
		logrus.WithError(err).Error("failed to build StorageAdapter")
		return nil, err
	}

	cmd := command.New()
	pdfConverter := pdfconv.New(htmltopdf.New(cmd))

	ci, err := certissuer.New(req.Issuer, req.ProcessId, req.Filename, storageAdapter, cmd, pdfConverter)
	if err != nil {
		logrus.WithError(err).Error("failed to build CertIssuer")
		return nil, err
	}

	logrus.Infof("Start issuing process: %s %s %s", req.Issuer, req.ProcessId, req.Filename)
	if err = ci.IssueCertificate(); err != nil {
		logrus.WithError(err).Error("failed cert_issuer.IssueCertificate")
		return nil, err
	}
	logrus.Infof("Finish issuing process: %s %s %s", req.Issuer, req.ProcessId, req.Filename)

	return &protocol.IssueBlockchainCertificateReply{}, nil
}
