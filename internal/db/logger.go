package db

import (
	"io"
	"log"

	"github.com/LassiHeikkila/taskey/internal/constants"
)

var logger = log.New(io.Discard, "[DB] ", constants.LogFlags)

func SetLoggerOutput(w io.Writer) {
	logger.SetOutput(w)
}
