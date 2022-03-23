package blobstore

type Service interface {
	Download(filename string) (interface{}, error)
	Upload(filename string, contentType string, filesize int64, data interface{}) error
	Delete(filename string) error
	List(prefix string) (interface{}, error)
}
