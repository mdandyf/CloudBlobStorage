package blobstore

type Config struct {
	AccountName   string `envconfig:"ACCOUNT_NAME" default:"test"`
	AccountKey    string `envconfig:"ACCOUNT_KEY" default:"123456"`
	ServiceURL    string `envconfig:"ACCOUNT_SERVICE_URL" default:"http://test"`
	ContainerName string `envconfig:"CONTAINER_NAME" default:"petrosea"` // in S3 is called BucketName
	Endpoint      string `envconfig:"ENDPOINT" default:"petrosea"`
	Region        string `envconfig:"REGION" default:"us-west-2"`
	SSL           bool   `envconfig:"SSL" default:"true"`
}
