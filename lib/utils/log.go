package utils

import (
	"log"
	"os"
	"runtime/debug"
	"strings"
)

// LogError prints error and exits if error is not nil
func LogError(err error) {
	if err != nil {
		log.Printf("%v\n\n", err)
		log.Println("Original error line:")

		stackLines := strings.Split(string(debug.Stack()), "\n")
		for iLine, line := range stackLines {
			if strings.Contains(line, "utils.LogError") {
				stackLines = stackLines[iLine+2 : iLine+4]
				break
			}
		}

		for _, line := range stackLines {
			log.Println(line)
		}
		os.Exit(1)
	}
}
