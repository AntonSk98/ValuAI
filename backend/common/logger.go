package common

import "log"

// LogError logs the error in a consistent format.
func LogError(err error) {
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
}

// LogInfo logs informational messages in a consistent format.
func LogInfo(msg string) {
	log.Printf("[INFO] %s", msg)
}
