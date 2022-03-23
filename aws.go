package blobstore

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func AwsStore(config *Config) Service {
	return &awsStore{
		config:  config,
		session: connectAwsBlobStorage(config),
		ctx:     context.Background(),
	}
}

type awsStore struct {
	config  *Config
	session *session.Session
	ctx     context.Context
}

func (a awsStore) Download(filename string) (interface{}, error) {
	downloader := s3manager.NewDownloader(a.session)

	var file *os.File
	_, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(a.config.ContainerName),
			Key:    aws.String(a.config.AccountKey),
		})

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (a awsStore) Upload(filename string, contentType string, filesize int64, data interface{}) error {
	// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
	// for more information on configuring part size, and concurrency.
	//
	// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
	uploader := s3manager.NewUploader(a.session)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(a.config.ContainerName),
		Key:    aws.String(filename),
		Body:   data.(*os.File),
	})

	if err != nil {
		return err
	}

	return nil
}

func (a awsStore) Delete(filename string) error {
	// Create S3 service client
	svc := s3.New(a.session)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(a.config.ContainerName), Key: aws.String(a.config.AccountKey)})
	if err != nil {
		return err
	}

	return nil
}

func (a awsStore) List(prefix string) (interface{}, error) {

	var results []*s3.Object

	svc := s3.New(a.session)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(a.config.ContainerName)})
	if err != nil {
		return nil, err
	}

	results = append(results, resp.Contents...)

	return results, nil
}

func connectAwsBlobStorage(config *Config) *session.Session {
	ses, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region)},
	)

	if err != nil {
		fmt.Println(err)
	}

	// Create S3 service client
	return ses
}
