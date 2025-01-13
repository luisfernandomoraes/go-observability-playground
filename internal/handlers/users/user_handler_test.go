package users

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TransformUsersSuite is the test suite for TransformUsersHandler.
type TransformUsersSuite struct {
	suite.Suite
	tests []struct {
		name           string
		method         string
		payload        string
		expectedStatus int
		expectedError  string
		validateResp   func(t *testing.T, resp TransformUsersResponse)
	}
}

// SetupSuite sets up the test cases.
func (suite *TransformUsersSuite) SetupSuite() {
	suite.tests = []struct {
		name           string
		method         string
		payload        string
		expectedStatus int
		expectedError  string
		validateResp   func(t *testing.T, resp TransformUsersResponse)
	}{
		{
			name:           "Valid Request",
			method:         http.MethodPost,
			payload:        `[{"name": "Alice", "age": 25}, {"name": "Bob", "age": 35}]`,
			expectedStatus: http.StatusOK,
			expectedError:  "",
			validateResp: func(t *testing.T, resp TransformUsersResponse) {
				assert.NotEmpty(t, resp.FilteredUsers, "FilteredUsers should not be empty")
				assert.Equal(t, "Bob", resp.FilteredUsers[0].Name, "FilteredUsers should contain Bob")
				assert.NotEmpty(t, resp.SortedUsers, "SortedUsers should not be empty")
				assert.NotNil(t, resp.GroupedUsers, "GroupedUsers should not be nil")
				assert.NotEmpty(t, resp.UpdatedUsers, "UpdatedUsers should not be empty")
				assert.GreaterOrEqual(t, resp.CountUsersAbove30, 1, "CountUsersAbove30 should be >= 1")
			},
		},
		{
			name:           "Invalid JSON Payload",
			method:         http.MethodPost,
			payload:        `{"name": "Alice", "age": "invalid"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Failed to decode JSON payload",
		},
		{
			name:           "Invalid Method",
			method:         http.MethodGet,
			payload:        "",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedError:  "Invalid request method.",
		},
		{
			name:           "Validation Error",
			method:         http.MethodPost,
			payload:        `[{"name": "", "age": 25}]`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed.",
		},
	}
}

// TestTransformUsersHandler runs all the test cases.
func (suite *TransformUsersSuite) TestTransformUsersHandler() {
	for _, tt := range suite.tests {
		suite.Run(tt.name, func() {
			// Prepare request
			req := httptest.NewRequest(tt.method, "/users", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			// Execute handler
			TransformUsersHandler(rr, req)

			// Assert status code
			suite.Equal(tt.expectedStatus, rr.Code, "Unexpected HTTP status")

			// Parse and validate response
			var resp TransformUsersResponse
			err := json.Unmarshal(rr.Body.Bytes(), &resp)
			if assert.NoError(suite.T(), err, "Response should be valid JSON") {
				// Validate error message
				if tt.expectedError != "" {
					assert.Contains(suite.T(), resp.ErrorMessage, tt.expectedError, "ErrorMessage should contain expected text")
				} else {
					assert.Empty(suite.T(), resp.ErrorMessage, "ErrorMessage should be empty on success")
				}

				// Additional response validation if applicable
				if tt.validateResp != nil {
					tt.validateResp(suite.T(), resp)
				}
			}
		})
	}
}

// TestTransformUsersSuite runs the test suite.
func TestTransformUsersSuite(t *testing.T) {
	suite.Run(t, new(TransformUsersSuite))
}
