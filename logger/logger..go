package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

// initiating custom logger
func init() {
	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	Log = slog.New(logHandler)
}
