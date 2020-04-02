package service

import (
	"context"

	"github.com/lastrust/issuing-service/config"
	"github.com/lastrust/issuing-service/config/di_container"
	"github.com/lastrust/issuing-service/domain/cert_issuer"
	"github.com/lastrust/issuing-service/protocol"
	"github.com/sirupsen/logrus"
)

type issuingService struct {
	conf *config.Config
}

// IssueBlockchainCertificate run the command of pkg/cert-issuer, returns an error if is not success
func (s issuingService) IssueBlockchainCertificate(ctx context.Context, req *protocol.IssueBlockchainCertificateRequest) (*protocol.IssueBlockchainCertificateReply, error) {
	storageAdapter, err := di_container.GetStorageAdapter(s.conf)
	if err != nil {
		logrus.WithError(err).Error("failed to build StorageAdapter")
		return nil, err
	}

	c, err := cert_issuer.New(req.Issuer, req.Filename, storageAdapter)
	if err != nil {
		logrus.WithError(err).Error("failed to build CertIssuer")
		return nil, err
	}

	err = c.IssueCertificate()
	if err != nil {
		logrus.WithError(err).Error("failed cert_issuer.IssueCertificate")
		return nil, err
	}

	return &protocol.IssueBlockchainCertificateReply{}, nil
}

// New issuingService constructor
func New(conf *config.Config) protocol.IssuingServiceServer {
	return &issuingService{conf}
}
