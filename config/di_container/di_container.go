package di_container

import (
	"errors"

	"github.com/lastrust/issuing-service/config"
	"github.com/lastrust/issuing-service/domain/certissuer"
	"github.com/lastrust/issuing-service/infra"
)

var errInvalidCloudService = errors.New("Invalid CLOUD_SERVICE")

func GetStorageAdapter(conf *config.Config) (certissuer.StorageAdapter, error) {
	switch conf.CloudService {
	case "GCP":
		return infra.NewGcsAdapter(conf.ProcessEnv)
	default:
		return nil, errInvalidCloudService
	}
}
