package utils

import "log"

// LogError prints error and exits if error is not nil.
func LogError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
