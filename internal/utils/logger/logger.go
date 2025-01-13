package logger

import (
	"io"
	"log"
	"log/slog"
	"os"
)

// SetupLogger sets up the application logger.
func SetupLogger(options ...string) {
	var output io.Writer = os.Stdout

	if len(options) > 0 && options[0] != "" {
		logFile := options[0]
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			output = os.Stderr
			log.Println("Failed to open log file:", err) // Use log.Printf to ensure it writes to stderr
		} else {
			output = file
		}
	}

	logger := slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{}))
	slog.SetDefault(logger)
}
