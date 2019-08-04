package awslib

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// CreateBucket .
func CreateBucket(c Storage) error {
	sess := Session(c)
	svc := s3.New(sess)

	_, err := svc.GetBucketVersioning(&s3.GetBucketVersioningInput{
		Bucket: aws.String(c.Bucket),
	})
	if err != nil {
		fmt.Println(err)

		_, err = svc.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(c.Bucket),
		})
		if err != nil {
			return err
		}
		fmt.Println("Bucket created:", c.Bucket)
	} else {
		fmt.Println("Bucket found:", c.Bucket)
	}

	return nil
}

// Upload .
func Upload(c Storage, bucket string, key string, file io.Reader) error {
	sess := Session(c)
	svc := s3.New(sess)
	uploader := s3manager.NewUploaderWithClient(svc)

	// Perform an upload.
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	return err
}

// Download .
func Download(c Storage, bucket string, key string) ([]byte, error) {
	sess := Session(c)
	svc := s3.New(sess)
	downloader := s3manager.NewDownloaderWithClient(svc)
	buf := aws.NewWriteAtBuffer([]byte{})

	// Perform a download.
	_, err := downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
