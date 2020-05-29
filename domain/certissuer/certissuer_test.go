package certissuer_test

import (
	"context"
	"os"
	"path/filepath"

	"github.com/golang/mock/gomock"
	"github.com/lastrust/issuing-service/domain/certissuer"
	"github.com/lastrust/issuing-service/mocks"
	"github.com/lastrust/issuing-service/utils/filesystem"
	"github.com/lastrust/issuing-service/utils/path"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Certissuer.IssueCertificate", func() {
	dataDir, _ := filepath.Abs("../../test/certissuer/data")
	path.SetDataDir(dataDir + "/")

	defer func() {
		path.SetDataDir("/storage/data/")
	}()

	Context("when indicate issuerId and processId of non existence configuration file", func() {
		ci := certissuer.New("", "", nil, nil, nil, nil)
		err := ci.IssueCertificate(context.Background())
		It("got an error ErrNoConfig", func() {
			Expect(err).To(Equal(certissuer.ErrNoConfig))
		})
	})

	Context("success case where IssueBlockchainCertificate and BulkCreate returns without error for empty certificates list", func() {
		ctrl := gomock.NewController(GinkgoT())
		defer ctrl.Finish()

		mockedCommand := mocks.NewMockCommand(ctrl)
		mockedCommand.EXPECT().IssueBlockchainCertificate(gomock.Any()).Return(nil)

		mockedCertRepo := mocks.NewMockRepository(ctrl)
		mockedCertRepo.EXPECT().BulkCreate(gomock.Any(), gomock.Any()).Return(nil)

		const (
			issuerId  = "test-issuer"
			processId = "test-process.IssueCertificate"
		)

		confPath := path.IssuerConfigPath(issuerId, processId)
		f, _ := os.OpenFile(confPath, os.O_CREATE, 0755)
		f.Close()

		ci := certissuer.New(issuerId, processId, nil, mockedCommand, nil, mockedCertRepo)
		err := ci.IssueCertificate(context.Background())
		defer os.RemoveAll(path.BlockcertsProcessDir(issuerId, processId))

		It("got an nil error", func() {
			Expect(err).To(BeNil())
		})
		It("issuer config file doesn't exists", func() {
			Expect(filesystem.FileExists(confPath)).To(Equal(false))
		})
	})
})
