package validations

import (
	"github.com/luisfernandomoraes/observability-go/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUsers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		input         []models.User
		expectedError string
	}{
		{
			name: "Valid Users",
			input: []models.User{
				{Name: "Alice", Age: 25},
				{Name: "Bob", Age: 35},
			},
			expectedError: "",
		},
		{
			name: "User With Empty Name",
			input: []models.User{
				{Name: "", Age: 25},
			},
			expectedError: "user at index 0 has an empty name",
		},
		{
			name: "User With Negative Age",
			input: []models.User{
				{Name: "Alice", Age: -5},
			},
			expectedError: "user at index 0 has an invalid age: -5",
		},
		{
			name: "User With Age Above 130",
			input: []models.User{
				{Name: "Bob", Age: 150},
			},
			expectedError: "user at index 0 has an invalid age: 150",
		},
		{
			name: "Multiple Users - First With Error",
			input: []models.User{
				{Name: "", Age: 25},
				{Name: "Bob", Age: 35},
			},
			expectedError: "user at index 0 has an empty name",
		},
		{
			name: "Multiple Users - Second With Error",
			input: []models.User{
				{Name: "Alice", Age: 25},
				{Name: "Bob", Age: 200},
			},
			expectedError: "user at index 1 has an invalid age: 200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Execute the function
			err := ValidateUsers(tt.input)

			// Assert expectations
			if tt.expectedError == "" {
				assert.NoError(t, err, "Expected no error for valid input")
			} else {
				assert.Error(t, err, "Expected an error for invalid input")
				assert.Contains(t, err.Error(), tt.expectedError, "Error message should match expected")
			}
		})
	}
}
