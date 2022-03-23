package blobstore

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func AzureStore(config *Config) Service {
	return &azureStore{
		config:    config,
		container: connectAzureBlobStorage(&config.AccountName, &config.AccountKey, &config.ServiceURL, &config.ContainerName),
		ctx:       context.Background(),
	}
}

type azureStore struct {
	config    *Config
	container azblob.ContainerURL
	ctx       context.Context
}

func (a azureStore) Download(filename string) (interface{}, error) {
	// Create a URL that references a to-be-created blob in your Azure Storage account's container.
	// This returns a BlockBlobURL object that wraps the blob's URL and a request pipeline (inherited from containerURL)
	blobURL := a.container.NewBlockBlobURL(filename) // Blob names can be mixed case

	// Download the blob's contents and verify that it worked correctly
	res, err := blobURL.Download(a.ctx, 0, 0, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a azureStore) Upload(filename string, contentType string, filesize int64, data interface{}) error {
	// Create a URL that references a to-be-created blob in your Azure Storage account's container.
	// This returns a BlockBlobURL object that wraps the blob's URL and a request pipeline (inherited from containerURL)
	blobURL := a.container.NewBlockBlobURL(filename) // Blob names can be mixed case

	// Upload the blob
	_, err := blobURL.Upload(a.ctx, data.(io.ReadSeeker), azblob.BlobHTTPHeaders{ContentType: contentType}, azblob.Metadata{}, azblob.BlobAccessConditions{}, azblob.DefaultAccessTier, nil, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (a azureStore) Delete(filename string) error {
	// Create a URL that references a to-be-created blob in your Azure Storage account's container.
	// This returns a BlockBlobURL object that wraps the blob's URL and a request pipeline (inherited from containerURL)
	blobURL := a.container.NewBlockBlobURL(filename) // Blob names can be mixed case

	// Delete the blob
	_, err := blobURL.Delete(a.ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		return err
	}

	return nil
}

func (a azureStore) List(prefix string) (interface{}, error) {
	var results [][]azblob.BlobItemInternal

	// List the blob(s) in our container; since a container may hold millions of blobs, this is done 1 segment at a time.
	for marker := (azblob.Marker{}); marker.NotDone(); { // The parens around Marker{} are required to avoid compiler error.
		// Get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := a.container.ListBlobsFlatSegment(a.ctx, marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			return nil, err
		}
		// IMPORTANT: ListBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker

		// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Segment.BlobItems {
			fmt.Print("Blob name: " + blobInfo.Name + "\n")
		}

		results = append(results, listBlob.Segment.BlobItems)
	}

	return results, nil
}

func connectAzureBlobStorage(accountName *string, accountKey *string, blobServiceURL *string, containerName *string) azblob.ContainerURL {

	// Use your Storage account's name and key to create a credential object; this is used to access your account.
	credential, _ := azblob.NewSharedKeyCredential(*accountName, *accountKey)

	// Create a request pipeline that is used to process HTTP(S) requests and responses. It requires
	// your account credentials. In more advanced scenarios, you can configure telemetry, retry policies,
	// logging, and other options. Also, you can configure multiple request pipelines for different scenarios.
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// From the Azure portal, get your Storage account blob service URL endpoint.
	// The URL typically looks like this:
	u, _ := url.Parse(fmt.Sprintf(*blobServiceURL, *accountName))

	// Create an ServiceURL object that wraps the service URL and a request pipeline.

	return azblob.NewServiceURL(*u, p).NewContainerURL(*containerName)
}
