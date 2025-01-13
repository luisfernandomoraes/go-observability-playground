package json

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

// DecodeJSONPayload decodes JSON from the request body into a given target struct.
func DecodeJSONPayload(r *http.Request, target interface{}) error {
	if r.Body == nil {
		slog.Warn("Empty request body", "url", r.URL.Path, "method", r.Method)
		return errors.New("empty request body")
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("Failed to close request body", "error", err, "url", r.URL.Path, "method", r.Method)
		}
	}(r.Body)

	// Check if the body is empty
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read request body", "error", err, "url", r.URL.Path, "method", r.Method)
		return errors.New("failed to read request body")
	}
	if len(bodyBytes) == 0 {
		slog.Warn("Empty request body", "url", r.URL.Path, "method", r.Method)
		return errors.New("empty request body")
	}

	// Decode JSON
	err = json.Unmarshal(bodyBytes, target)
	if err != nil {
		slog.Error("Failed to decode JSON payload", "error", err, "url", r.URL.Path, "method", r.Method)
		return errors.New("failed to decode JSON payload")
	}

	slog.Info("Successfully decoded JSON payload", "url", r.URL.Path, "method", r.Method)
	return nil
}

// WriteJSONResponse writes a JSON response with a given status code and logger any errors during encoding.
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(data)
	if err != nil {
		slog.Error("Failed to write JSON response",
			"error", err,
			"statusCode", statusCode,
			"data", data,
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	slog.Info("Successfully wrote JSON response",
		"statusCode", statusCode,
		"data", data,
	)
}
