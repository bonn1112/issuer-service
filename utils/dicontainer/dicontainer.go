package dicontainer

import (
	"errors"

	"github.com/lastrust/issuing-service/domain/certissuer"
	"github.com/lastrust/issuing-service/infra/adapter/gcs"
	"github.com/lastrust/issuing-service/infra/adapter/s3"
)

var errInvalidCloudService = errors.New("Invalid CLOUD_SERVICE")

func GetStorageAdapter(cloudService, processEnv string) (certissuer.StorageAdapter, error) {
	switch cloudService {
	case "gcp":
		return gcs.New(processEnv)
	case "aws":
		return s3.New(processEnv)
	default:
		return nil, errInvalidCloudService
	}
}
