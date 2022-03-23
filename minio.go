package blobstore

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func MinioStore(config *Config) Service {
	return &minioStore{
		config: config,
		client: connectMinioBlobStorage(config),
		ctx:    context.Background(),
	}
}

type minioStore struct {
	config *Config
	client *minio.Client
	ctx    context.Context
}

func (m minioStore) Download(filename string) (interface{}, error) {
	object, err := m.client.GetObject(context.Background(), m.config.AccountName, filename, minio.GetObjectOptions{
		ServerSideEncryption: nil,
	})
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (m minioStore) Upload(filename string, contentType string, filesize int64, data interface{}) error {
	_, err := m.client.PutObject(context.Background(), m.config.AccountName, filename, data.(io.Reader), filesize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}

	return nil
}

func (m minioStore) Delete(filename string) error {
	return m.client.RemoveObject(m.ctx, m.config.AccountName, filename, minio.RemoveObjectOptions{
		ForceDelete: true,
	})
}

func (m minioStore) List(prefix string) (interface{}, error) {
	ctx, cancel := context.WithCancel(m.ctx)

	defer cancel()

	objects := m.client.ListObjects(ctx, m.config.AccountName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	var results []*minio.ObjectInfo

	for object := range objects {
		if object.Err != nil {
			return nil, object.Err
		}

		objectInfo, _ := m.client.StatObject(ctx, m.config.AccountName, object.Key, minio.GetObjectOptions{})
		results = append(results, &objectInfo)
	}

	return results, nil
}

func connectMinioBlobStorage(config *Config) *minio.Client {
	client, _ := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccountKey, config.AccountSecret, ""),
		Secure: config.SSL,
	})

	return client
}