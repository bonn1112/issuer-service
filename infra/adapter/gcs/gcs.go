package gcs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/lastrust/issuing-service/utils/path"
)

var errInvalidProcessEnv = errors.New("Invalid PROCESS_ENV")

type gcsAdapter struct {
	bucket string
}

func New(processEnv string) (*gcsAdapter, error) {
	var bucket string
	switch processEnv {
	case "dev":
		bucket = "lst-issuer-dev"
	case "stg":
		bucket = "lst-issuer-stg"
	case "prd":
		bucket = "lst-issuer-prd"
	default:
		return nil, errInvalidProcessEnv
	}

	return &gcsAdapter{bucket}, nil
}

func (s *gcsAdapter) StoreCertificate(filepath, issuerId, filename string) (err error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	pathInGcs := path.CertsPathInGCS(issuerId, filename)
	w := client.Bucket(s.bucket).Object(pathInGcs).NewWriter(ctx)
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}

	if _, err = io.Copy(w, f); err != nil {
		return err
	}

	return w.Close()
}

func (adapter *gcsAdapter) StorePdf(_, _, _ string) error {
	return fmt.Errorf("gcsAdapter.SrorePdf is not implemented")
}
