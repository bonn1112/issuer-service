package s3

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type s3Adapter struct {
	bucket string
}

func New(processEnv string) (*s3Adapter, error) {
	var bucket string
	switch processEnv {
	case "prd":
		bucket = "issued-cloudcerts-prd"
	case "stg", "dev":
		fallthrough
	default:
		bucket = "issued-cloudcerts-stg"
	}

	return &s3Adapter{bucket}, nil
}

func (adapter *s3Adapter) StoreCertificate(filepath, issuerId, filename string) error {
	pathInS3 := fmt.Sprintf("%s/%s", issuerId, filename)
	return adapter.upload(filepath, pathInS3)
}

func (adapter *s3Adapter) StorePdf(filepath, issuerId, filenameWithoutExt string) error {
	pathInS3 := fmt.Sprintf("pdf/%s/%s.pdf", issuerId, filenameWithoutExt)
	return adapter.upload(filepath, pathInS3)
}

func (adapter *s3Adapter) StorePasswordRecords(issuerId, processId string, records [][]string) error {
	pathInS3 := fmt.Sprintf("password_records/%s/%s.csv", issuerId, processId)

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.WriteAll(records); err != nil {
		return err
	}
	w.Flush()

	return adapter.uploadFromReader(bytes.NewReader(buf.Bytes()), pathInS3)
}

func (adapter *s3Adapter) upload(filepath, key string) error {
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

	_, err = uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(adapter.bucket),
		Key:    aws.String(key),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	return nil
}

func (adapter *s3Adapter) uploadFromReader(reader io.Reader, key string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	uploader := s3manager.NewUploader(sess)

	_, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(adapter.bucket),
		Key:    aws.String(key),
		Body:   reader,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	return nil
}
