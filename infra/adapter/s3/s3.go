package s3

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/lastrust/issuing-service/utils/path"
)

var errInvalidProcessEnv = errors.New("Invalid PROCESS_ENV")

type s3Adapter struct {
	bucket string
}

func New(processEnv string) (*s3Adapter, error) {
	var bucket string
	switch processEnv {
	case "dev":
		bucket = "issued-cloudcerts-stg"
	case "stg":
		bucket = "issued-cloudcerts-stg"
	case "prd":
		bucket = "issued-cloudcerts-prd"
	default:
		return nil, errInvalidProcessEnv
	}

	return &s3Adapter{bucket}, nil
}

func (s *s3Adapter) StoreCerts(filepath string, issuer string, filename string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	uploader := s3manager.NewUploader(sess)

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	pathInS3 := path.CertsPathInS3(issuer, filename)

	_, err = uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(pathInS3),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	return nil
}
