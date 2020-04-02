package di_container

import (
	"errors"

	"github.com/lastrust/issuing-service/config"
	"github.com/lastrust/issuing-service/domain/cert_issuer"
	"github.com/lastrust/issuing-service/infra"
)

var errInvalidCloudService = errors.New("Invalid CLOUD_SERVICE")

func GetStorageAdapter(conf *config.Config) (cert_issuer.StorageAdapter, error) {
	switch conf.CloudService {
	case "GCP":
		return infra.NewGcsAdapter(conf.ProcessEnv)
	default:
		return nil, errInvalidCloudService
	}
}
