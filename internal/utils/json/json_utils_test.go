package json

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecodeJSONPayload(t *testing.T) {
	tests := []struct {
		name          string
		requestBody   string
		expectedError error
		expectedData  map[string]interface{}
	}{
		{
			name:          "Valid JSON payload",
			requestBody:   `{"name": "Alice", "age": 30}`,
			expectedError: nil,
			expectedData: map[string]interface{}{
				"name": "Alice",
				"age":  float64(30), // json.Decoder converts numbers to float64 by default
			},
		},
		{
			name:          "Empty request body",
			requestBody:   "",
			expectedError: errors.New("empty request body"),
			expectedData:  nil,
		},
		{
			name:          "Invalid JSON payload",
			requestBody:   `{"name": "Alice", "age":}`,
			expectedError: errors.New("failed to decode JSON payload"),
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(tt.requestBody))

			var target map[string]interface{}
			err := DecodeJSONPayload(req, &target)

			// Check for expected errors
			if tt.expectedError != nil {
				require.EqualError(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedData, target, "Decoded data does not match expected")
			}
		})
	}
}

func TestWriteJSONResponse(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseData   interface{}
		expectedBody   string
		expectedStatus int
	}{
		{
			name:           "Valid JSON response",
			statusCode:     http.StatusOK,
			responseData:   map[string]interface{}{"message": "Success"},
			expectedBody:   `{"message":"Success"}` + "\n", // json.Encoder adds a newline
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Internal Server Error on JSON encoding failure",
			statusCode:     http.StatusInternalServerError,
			responseData:   func() {}, // Invalid type for JSON encoding
			expectedBody:   "Failed to encode response\n",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			// Execute WriteJSONResponse
			WriteJSONResponse(rr, tt.statusCode, tt.responseData)

			// Assert HTTP status code
			require.Equal(t, tt.expectedStatus, rr.Code)

			// Assert response body
			require.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
