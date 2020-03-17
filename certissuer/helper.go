package certissuer

import (
	"context"

	"github.com/lastrust/issuing-service/config"

	"cloud.google.com/go/storage"
	"io/ioutil"

)

func getGCSClient() (client *Client, bucket string, err error) {
	var bucket string

	conf, err := config.Env()
	if err != nil {
		return nil, nil, err
	}
	if conf.ProcessEnv == "dev" || conf.ProcessEnv == "stg" {
		bucket = "lastrust-stg"
	} else if conf.ProcessEnv == "prd" {
		bucket = "lastrust-prd"
	} else {
		return nil, nil, errors.New("Invalid PROCESS_ENV")
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != {
		return nil, nil, err
	}

	return client, bucket, nil
}

func storeGCS(filepath string) (err error) {
	client, bucket, err = getGCSClient()
	if err != nil {
		return nil, err
	}
	w, err := client.Bucket(bucket).Object(object).NewWriter(ctx)
	if err != nil {
			return nil, err
	}
	defer w.Close()

	f, err := os.Open(filepath)
	_, err = io.Copy(wc, f)
	if err != nil {
		return err
	}
	err := wc.Close(
	if err != nil {
		return err
	}
	return
}