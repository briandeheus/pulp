package helpers

import "os"

func GetEnvVarString(key string, def string) string {

	val := os.Getenv(key)

	if val == "" {
		return def
	}

	return val

}