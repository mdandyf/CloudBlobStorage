package blobstore

type Config struct {
	AccountName   string
	AccountKey    string
	AccountSecret string
	ServiceURL    string
	ContainerName string
	Endpoint      string
	SSL           bool
}
