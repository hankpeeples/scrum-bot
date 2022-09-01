// Package utils contains utility functions needed throughout the app
package utils

import (
	"os"

	logger "github.com/withmandala/go-log"
)

// Prefix is the bot command character prefix
const Prefix string = "?"

var log = NewLogger()

// NewLogger returns a new instance of a logger
func NewLogger() *logger.Logger {
	return logger.New(os.Stdout).WithColor()
}
