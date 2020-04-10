package dicontainer

import (
	"errors"

	"github.com/lastrust/issuing-service/domain/certissuer"
	"github.com/lastrust/issuing-service/infra"
)

var errInvalidCloudService = errors.New("Invalid CLOUD_SERVICE")

func GetStorageAdapter(cloudService, processEnv string) (certissuer.StorageAdapter, error) {
	switch cloudService {
	case "gcp":
		return infra.NewGcsAdapter(processEnv)
	default:
		return nil, errInvalidCloudService
	}
}
