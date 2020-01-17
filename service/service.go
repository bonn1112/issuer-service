package service

import (
	"context"

	"github.com/lastrust/issuing-service/certissuer"
	"github.com/lastrust/issuing-service/protocol"
)

type issuingService struct{}

// IssueBlockchainCertificate run the command of pkg/cert-issuer, returns an error if is not success
func (s issuingService) IssueBlockchainCertificate(ctx context.Context, req *protocol.IssueBlockchainCertificateRequest) (*protocol.IssueBlockchainCertificateReply, error) {
	cli := certissuer.New(req.Issuer, req.Filename)

	err := cli.IssueCertificate()
	if err != nil {
		return nil, err
	}

	return &protocol.IssueBlockchainCertificateReply{}, nil
}

// New issuingService constructor
func New() protocol.IssuingServiceServer {
	return &issuingService{}
}
