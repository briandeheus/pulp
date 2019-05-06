package localFiles

import (
	"github.com/google/uuid"
	"log"
	"os"
	"path"
)

var tmpLocation = path.Join("/", "tmp", uuid.New().String())
var cacheMap map[string]bool

func init () {

	log.Printf("Creating temporary location for files in %s", tmpLocation)

	if err := os.MkdirAll(tmpLocation, os.ModePerm); err != nil {
		log.Fatalf("Failed to create location in %s", tmpLocation)
		panic(err)
	}

	cacheMap = make(map[string]bool)

}

func getExistsLocal(filePath *string) bool {

	val, exists := cacheMap[*filePath]

	if exists == true {
		return val
	}

	return false

}

func SetExistsLocal(filePath *string, exists bool) {

	cacheMap[*filePath] = exists

}

func GetLocalPath(filePath *string) string {
	return path.Join(tmpLocation, *filePath)
}

func Exists(filePath *string) bool {

	fullPath := GetLocalPath(filePath)

	if getExistsLocal(filePath) == true {

		log.Printf("file=%s exists=true source=cache", fullPath)
		return true

	} else if _, err := os.Stat(fullPath); os.IsNotExist(err) {

		log.Printf("file=%s exists=no source=%v", fullPath, err)
		return false

	} else {

		SetExistsLocal(filePath, true)
		log.Printf("file=%s exists=true source=disk", fullPath)
		return true

	}

}

func Delete(filePath *string) error {

	fullPath := GetLocalPath(filePath)

	// Make sure we check whether the file exists or not before deleting.
	// If the file does not exist, we don't really care.
	if Exists(filePath) == false {
		return nil
	}

	// Then remove.
	return os.Remove(fullPath)

}