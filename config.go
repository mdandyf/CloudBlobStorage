package blobstore

type Config struct {
	AccountName   string `envconfig:"ACCOUNT_NAME" default:"test"`
	AccountKey    string `envconfig:"ACCOUNT_KEY" default:"123456"`
	ServiceURL    string `envconfig:"ACCOUNT_SERVICE_URL" default:"http://test"`
	ContainerName string `envconfig:"CONTAINER_NAME" default:"petrosea"`
	Endpoint      string `envconfig:"ENDPOINT" default:"petrosea"`
	SSL           bool   `envconfig:"SSL" default:"true"`
}
