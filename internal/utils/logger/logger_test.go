package logger

import (
	"bytes"
	"io"
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupLogger(t *testing.T) {
	t.Run("Default to Stdout", func(t *testing.T) {
		// Create a pipe to capture os.Stdout
		r, w, err := os.Pipe()
		require.NoError(t, err)

		// Replace os.Stdout with the write end of the pipe
		stdout := os.Stdout
		defer func() { os.Stdout = stdout }()
		os.Stdout = w

		// Run SetupLogger
		SetupLogger()

		// Log a test message
		slog.Info("Test log message to stdout")

		// Close the write end of the pipe and capture output
		w.Close()
		var buf bytes.Buffer
		_, err = io.Copy(&buf, r)
		require.NoError(t, err)

		// Check output
		require.Contains(t, buf.String(), "Test log message to stdout")
	})

	t.Run("Log to File", func(t *testing.T) {
		// Create a temporary file
		tempFile, err := os.CreateTemp("", "test-log-*.log")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())

		// Run SetupLogger with the temporary file
		SetupLogger(tempFile.Name())

		// Log a test message
		slog.Info("Test log message to file")

		// Read file content
		fileContent, err := os.ReadFile(tempFile.Name())
		require.NoError(t, err)

		// Check output
		assert.Contains(t, string(fileContent), "Test log message to file")
	})

	t.Run("Fallback to Stderr on File Error", func(t *testing.T) {
		// Create a pipe to capture os.Stderr
		r, w, err := os.Pipe()
		require.NoError(t, err)

		// Create a copy of os.Stderr to a local variable
		stderr := os.Stderr
		// Restore os.Stderr after the test completes
		defer func() { os.Stderr = stderr }()
		// Replace os.Stderr with the write end of the pipe
		os.Stderr = w

		// Redirect standard logger output to the pipe
		log.SetOutput(w)

		// Run SetupLogger with an invalid file path
		SetupLogger("/invalid/path/to/logfile.log")

		// Log a test message
		slog.Info("Test log message to stderr")

		// Close the write end of the pipe and capture output
		w.Close()
		var buf bytes.Buffer
		_, err = io.Copy(&buf, r)
		require.NoError(t, err)

		// Check output
		assert.Contains(t, buf.String(), "Failed to open log file")
		assert.Contains(t, buf.String(), "Test log message to stderr")
	})
}
