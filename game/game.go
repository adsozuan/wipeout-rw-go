package wipeout

import (
	"log"
	"os"
)

// Logger is a package-level logger
var Logger *log.Logger

func init() {
	// Create a logger that writes to standard error with a timestamp
	Logger = log.New(os.Stderr, "wipeout|", log.Ldate|log.Ltime)
}
