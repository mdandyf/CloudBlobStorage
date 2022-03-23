package blobstore

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func GcpStore(config *Config) Service {
	return &gcpStore{
		config: config,
		client: connectGcpBlobStorage(config),
		ctx:    context.Background(),
	}
}

type gcpStore struct {
	config *Config
	client *storage.Client
	ctx    context.Context
}

func (g gcpStore) Download(filename string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(g.ctx, time.Second*100)
	defer cancel()

	f, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("os.Create: %v", err)
	}

	rc, err := g.client.Bucket(g.config.ContainerName).Object(filename).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", filename, err)
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return nil, fmt.Errorf("io.Copy: %v", err)
	}

	if err = f.Close(); err != nil {
		return nil, fmt.Errorf("f.Close: %v", err)
	}

	return rc, nil
}

func (g gcpStore) Upload(filename string, contentType string, filesize int64, data interface{}) error {
	ctx, cancel := context.WithTimeout(g.ctx, time.Second*100)
	defer cancel()

	o := g.client.Bucket(g.config.ContainerName).Object(filename)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// For an object that does not yet exist, set the DoesNotExist precondition.
	o = o.If(storage.Conditions{DoesNotExist: true})
	// If the live object already exists in your bucket, set instead a
	// generation-match precondition using the live object's generation number.
	// attrs, err := o.Attrs(ctx)
	// if err != nil {
	//      return fmt.Errorf("object.Attrs: %v", err)
	// }
	// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	// Upload an object with storage.Writer.

	wc := o.NewWriter(ctx)
	if _, err := io.Copy(wc, data.(*os.File)); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

func (g gcpStore) Delete(filename string) error {
	ctx, cancel := context.WithTimeout(g.ctx, time.Second*10)
	defer cancel()

	o := g.client.Bucket(g.config.ContainerName).Object(filename)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %v", err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %v", filename, err)
	}

	return nil
}

func (g gcpStore) List(prefix string) (interface{}, error) {
	var results []*storage.ObjectAttrs

	// Prefixes and delimiters can be used to emulate directory listings.
	// Prefixes can be used to filter objects starting with prefix.
	// The delimiter argument can be used to restrict the results to only the
	// objects in the given "directory". Without the delimiter, the entire tree
	// under the prefix is returned.
	//
	// For example, given these blobs:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// If you just specify prefix="a/", you'll get back:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// However, if you specify prefix="a/" and delim="/", you'll get back:
	//   /a/1.txt
	ctx, cancel := context.WithTimeout(g.ctx, time.Second*10)
	defer cancel()

	it := g.client.Bucket(g.config.ContainerName).Objects(ctx, &storage.Query{
		Prefix:    prefix,
		Delimiter: "",
	})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Bucket(%q).Objects(): %v", g.config.ContainerName, err)
		}

		results = append(results, attrs)
	}

	return results, nil
}

func connectGcpBlobStorage(config *Config) *storage.Client {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	defer client.Close()

	return client
}
