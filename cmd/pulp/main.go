package main

import (
	"github.com/briandeheus/pulp/internal/pkg/helpers"
	"github.com/briandeheus/pulp/internal/pkg/localFiles"
	"github.com/briandeheus/pulp/internal/pkg/remoteFiles"
	"html"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var indexDefault string

func init () {

	indexDefault = helpers.GetEnvVarString("PULP_INDEX", "index.html")

}

func hydrateHeaders(w http.ResponseWriter) {

	// Indicate to consumers that they're talking to a Pulp server
	w.Header().Set("Server", "Pulp")

}

func handleDelete(w http.ResponseWriter, filepath *string) {

	log.Printf("Deleting file=%s", *filepath)

	if err := localFiles.Delete(filepath); err != nil {

		// If there is an error we log, return 500.
		log.Printf("Failed to delete local file error=%s", *filepath)
		w.WriteHeader(500)

	} else {

		// Otherwise we update the cache, return 200.
		localFiles.SetExistsLocal(filepath, false)
		w.WriteHeader(200)

	}

}

func handleLocalFileExists(w http.ResponseWriter, r *http.Request, path *string) {

	w.Header().Set("X-Cache", "Hit")
	http.ServeFile(w, r, localFiles.GetLocalPath(path))
	return
}

func handleRemoteFileExists(w http.ResponseWriter, r *http.Request, path * string) {

	localPath := localFiles.GetLocalPath(path)
	pathWithoutFileName, _ := filepath.Split(localPath)

	// Create the directory where the file should live.
	// This is useful for nested directories.
	err := os.MkdirAll(pathWithoutFileName, os.ModePerm)

	if err != nil {

		// Something went bad, return 500.
		log.Panicf("Failed to create local directory: %v", err)
		w.WriteHeader(500)
		return

	}

	fileHandler, err := os.Create(localPath)

	if err != nil {

		// There is a chance we can't create a file handler, return 500.
		log.Printf("Failed to create filehandler: %s", localPath)
		w.WriteHeader(500)
		return

	}

	if err := remoteFiles.DownloadToFile(fileHandler, path); err != nil {

		// So much that can go wrong. As always, return 500 if it does.
		log.Printf("Failed to download file: %s", *path)
		w.WriteHeader(500)
		return

	}

	// Set the header indicating that we missed cache. Useful for debugging.
	w.Header().Set("X-Cache", "Miss")
	http.ServeFile(w, r, localPath)

}

func handleAll(w http.ResponseWriter, r *http.Request, path *string) {

	// If the path is empty it means we're requesting the index.
	if *path == "" {
		*path = indexDefault
	}

	// Check if the path exists locally first.
	// Will hit disk only if the file exists, and it hasn't been requested before.
	if localFiles.Exists(path) == true {
		handleLocalFileExists(w, r, path)
		return
	}

	// Check if the remote file exists.
	// If it doesn't, yeet back a 404.
	if remoteFiles.Exists(path) == false {
		w.WriteHeader(404)
		return
	}

	// The file exists remotely, time to fetch and serve.
	handleRemoteFileExists(w, r, path)

}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// We trim HTML entities and also the first leading slash.
		path := html.EscapeString(r.URL.Path)
		path = strings.TrimLeft(path, "/")

		hydrateHeaders(w)

		// If we encounter a delete it means we want to remove an entry from cache
		if r.Method == "DELETE" {
			handleDelete(w, &path)
		} else {
			handleAll(w, r, &path)
		}

	})

	address := helpers.GetEnvVarString("PULP_ADDRESS", "0.0.0.0:8000")

	log.Fatal(http.ListenAndServe(address, nil))

}
