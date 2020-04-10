package env

import "os"

func GetDefault(key, defaults string) string {
	if v := os.Getenv(key); v != "" {
		return v
	} else {
		return defaults
	}
}
