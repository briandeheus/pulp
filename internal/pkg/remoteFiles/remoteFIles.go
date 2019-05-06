package remoteFiles

import (
	"cloud.google.com/go/storage"
	"context"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var client storage.Client
var prefix string
var bucket string

func init() {

	prefix = os.Getenv("PULP_PREFIX")
	bucket = os.Getenv("PULP_BUCKET")

	if prefix == "" {
		log.Fatalf("PULP_PREFIX is not set")
	}

	if bucket == "" {
		log.Fatalf("PULP_BUCKET is not set")
	}

	ctx := context.Background()
	c, err := storage.NewClient(ctx)

	if err != nil {
		log.Panicf("Failed to initialize Google Storage client: %v", err)
	}

	client = *c

}

func getFullPath(filePath *string) string {

	return path.Join(prefix, *filePath)

}

func DownloadToFile(file *os.File, filePath *string) error {

	// Get the bucket + the full path of the object on GCS.
	// Full path = prefix + requested path.
	bkt := client.Bucket(bucket)
	fullPath := getFullPath(filePath)

	// Open the file
	reader, err := bkt.Object(fullPath).NewReader(context.Background())

	if err != nil {
		return err
	}

	// S U C C all the bytes in
	fileBytes, err := ioutil.ReadAll(reader)

	if err != nil {
		return err
	}

	// Write the file to disk.
	if _, err := file.Write(fileBytes); err != nil {
		return err
	}

	// Return with nil to indicate the file has been written.
	return nil

}


func Exists(filePath *string) bool {

	exists := true

	// Get the bucket and obj on GCS.
	// Full path = prefix + requested path.
	fullPath := getFullPath(filePath)
	bkt := client.Bucket(bucket)
	obj := bkt.Object(fullPath)

	// Object reader needs a context.
	readerCtx := context.Background()
	_, err := obj.NewReader(readerCtx)

	if err != nil {

		// ErrObjectNotExists is expected, everything else we log.
		// We consciously don't panic to keep the server up and running.
		if err != storage.ErrObjectNotExist {
			log.Printf("unable to read data from file=%s error=%v", *filePath, err)
		}

		exists = false

	}

	return exists

}
