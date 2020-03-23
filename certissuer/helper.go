package certissuer

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"github.com/lastrust/issuing-service/config"

	"cloud.google.com/go/storage"
)

func getBucket() (string, error) {
	conf, err := config.Env()
	if err != nil {
		return "", err
	}

	switch conf.ProcessEnv {
	case "dev":
		return "lst-issuer-dev", nil
	case "stg":
		return "lst-issuer-stg", nil
	case "prd":
		return "lst-issuer-prd", nil
	default:
		return "", errors.New("Invalid PROCESS_ENV")
	}
}

func (i *certIssuer) storeGCS(filepath string) (err error) {
	bucket, err := getBucket()
	if err != nil {
		return err
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	w := client.Bucket(bucket).Object(i.certsPathInGCS()).NewWriter(ctx)
	f, err := os.Open(filepath)
	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return
}
