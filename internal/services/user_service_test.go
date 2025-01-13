package services

import (
	"github.com/luisfernandomoraes/observability-go/internal/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestFilterUsersAboveAge(t *testing.T) {
	t.Parallel()
	// Define test cases
	tests := []struct {
		name           string
		users          []models.User
		age            int
		expectedResult []models.User
	}{
		{
			name: "All users above age",
			users: []models.User{
				{Name: "Luis", Age: 35},
				{Name: "Alice", Age: 40},
			},
			age: 30,
			expectedResult: []models.User{
				{Name: "Luis", Age: 35},
				{Name: "Alice", Age: 40},
			},
		},
		{
			name: "No users above age",
			users: []models.User{
				{Name: "Luis", Age: 20},
				{Name: "Alice", Age: 19},
			},
			age:            30,
			expectedResult: []models.User{},
		},
		{
			name: "Some users above age",
			users: []models.User{
				{Name: "Luis", Age: 25},
				{Name: "Alice", Age: 35},
				{Name: "Bob", Age: 30},
			},
			age: 30,
			expectedResult: []models.User{
				{Name: "Alice", Age: 35},
			},
		},
		{
			name:           "Empty user list",
			users:          []models.User{},
			age:            30,
			expectedResult: []models.User{},
		},
		{
			name: "Edge case: users with age equal to filter",
			users: []models.User{
				{Name: "Luis", Age: 30},
				{Name: "Alice", Age: 29},
				{Name: "Bob", Age: 31},
			},
			age: 30,
			expectedResult: []models.User{
				{Name: "Bob", Age: 31},
			},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := FilterUsersAboveAge(tt.users, tt.age)
			require.Equal(t, tt.expectedResult, result, "Filtered users do not match expected result")
		})
	}
}

func TestSortUsersByName(t *testing.T) {
	t.Parallel()
	// Define test cases
	tests := []struct {
		name           string
		users          []models.User
		expectedResult []models.User
	}{
		{
			name: "Sort users by name in alphabetical order",
			users: []models.User{
				{Name: "Luis", Age: 25},
				{Name: "Alice", Age: 30},
				{Name: "Bob", Age: 20},
			},
			expectedResult: []models.User{
				{Name: "Alice", Age: 30},
				{Name: "Bob", Age: 20},
				{Name: "Luis", Age: 25},
			},
		},
		{
			name: "Case sensitivity in sorting",
			users: []models.User{
				{Name: "alice", Age: 30},
				{Name: "Bob", Age: 20},
				{Name: "Luis", Age: 25},
			},
			expectedResult: []models.User{
				{Name: "Bob", Age: 20},
				{Name: "Luis", Age: 25},
				{Name: "alice", Age: 30},
			},
		},
		{
			name:           "Empty user list remains empty",
			users:          []models.User{},
			expectedResult: []models.User{},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := SortUsersByName(tt.users)
			require.Equal(t, tt.expectedResult, result, "Sorted users do not match expected result")
		})
	}
}

func TestGroupUsersByAge(t *testing.T) {
	t.Parallel()
	// Define test cases
	tests := []struct {
		name           string
		users          []models.User
		expectedGroups map[string][]models.User
	}{
		{
			name: "Users grouped in all ranges",
			users: []models.User{
				{Name: "Luis", Age: 19},
				{Name: "Alice", Age: 25},
				{Name: "Bob", Age: 35},
				{Name: "Charlie", Age: 45},
			},
			expectedGroups: map[string][]models.User{
				"<20":   {{Name: "Luis", Age: 19}},
				"20-29": {{Name: "Alice", Age: 25}},
				"30-39": {{Name: "Bob", Age: 35}},
				"40+":   {{Name: "Charlie", Age: 45}},
			},
		},
		{
			name:  "No users",
			users: []models.User{},
			expectedGroups: map[string][]models.User{
				"<20":   {},
				"20-29": {},
				"30-39": {},
				"40+":   {},
			},
		},
		{
			name: "All users in one range",
			users: []models.User{
				{Name: "Luis", Age: 15},
				{Name: "Alice", Age: 18},
				{Name: "Bob", Age: 19},
			},
			expectedGroups: map[string][]models.User{
				"<20":   {{Name: "Luis", Age: 15}, {Name: "Alice", Age: 18}, {Name: "Bob", Age: 19}},
				"20-29": {},
				"30-39": {},
				"40+":   {},
			},
		},
		{
			name: "Ages at range boundaries",
			users: []models.User{
				{Name: "Luis", Age: 20},
				{Name: "Alice", Age: 30},
				{Name: "Bob", Age: 40},
			},
			expectedGroups: map[string][]models.User{
				"<20":   {},
				"20-29": {{Name: "Luis", Age: 20}},
				"30-39": {{Name: "Alice", Age: 30}},
				"40+":   {{Name: "Bob", Age: 40}},
			},
		},
		{
			name: "Duplicate users across ranges",
			users: []models.User{
				{Name: "Luis", Age: 19},
				{Name: "Luis", Age: 25},
				{Name: "Luis", Age: 35},
				{Name: "Luis", Age: 45},
			},
			expectedGroups: map[string][]models.User{
				"<20":   {{Name: "Luis", Age: 19}},
				"20-29": {{Name: "Luis", Age: 25}},
				"30-39": {{Name: "Luis", Age: 35}},
				"40+":   {{Name: "Luis", Age: 45}},
			},
		},
	}

	// Execute the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := GroupUsersByAge(tt.users)
			for key, expectedGroup := range tt.expectedGroups {
				require.ElementsMatch(t, expectedGroup, result[key], "Group %s does not match", key)
			}
		})
	}
}

func TestUpdateUsersAge(t *testing.T) {
	t.Parallel()
	// Define test cases
	tests := []struct {
		name          string
		users         []models.User
		increment     int
		expectedAges  []int
		sleepDuration time.Duration
	}{
		{
			name: "Increment age by 1",
			users: []models.User{
				{Name: "Luis", Age: 20},
				{Name: "Alice", Age: 25},
			},
			increment:     1,
			expectedAges:  []int{21, 26},
			sleepDuration: time.Second, // Match the `time.Sleep` in the function
		},
		{
			name: "Increment age by 5",
			users: []models.User{
				{Name: "Luis", Age: 30},
				{Name: "Alice", Age: 35},
			},
			increment:     5,
			expectedAges:  []int{35, 40},
			sleepDuration: time.Second,
		},
		{
			name: "No increment",
			users: []models.User{
				{Name: "Luis", Age: 40},
				{Name: "Alice", Age: 45},
			},
			increment:     0,
			expectedAges:  []int{40, 45},
			sleepDuration: time.Second,
		},
	}

	// Run test cases
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Execute the function
			UpdateUsersAge(tt.users, tt.increment)

			// Verify the ages of the users
			actualAges := make([]int, len(tt.users))
			for i, user := range tt.users {
				actualAges[i] = user.Age
			}
			require.Equal(t, tt.expectedAges, actualAges, "User ages do not match expected values")
		})
	}
}

func TestCountUsersAboveAge(t *testing.T) {
	t.Parallel()
	// Define test cases
	tests := []struct {
		name          string
		users         []models.User
		age           int
		expectedCount int
	}{
		{
			name: "All users above age",
			users: []models.User{
				{Name: "Luis", Age: 35},
				{Name: "Alice", Age: 40},
			},
			age:           30,
			expectedCount: 2,
		},
		{
			name: "No users above age",
			users: []models.User{
				{Name: "Luis", Age: 20},
				{Name: "Alice", Age: 19},
			},
			age:           30,
			expectedCount: 0,
		},
		{
			name: "Some users above age",
			users: []models.User{
				{Name: "Luis", Age: 25},
				{Name: "Alice", Age: 35},
				{Name: "Bob", Age: 30},
			},
			age:           30,
			expectedCount: 2,
		},
		{
			name: "Boundary age condition",
			users: []models.User{
				{Name: "Luis", Age: 30},
				{Name: "Alice", Age: 29},
				{Name: "Bob", Age: 31},
			},
			age:           30,
			expectedCount: 2,
		},
		{
			name:          "Empty user list",
			users:         []models.User{},
			age:           30,
			expectedCount: 0,
		},
		{
			name: "All users at boundary age",
			users: []models.User{
				{Name: "Luis", Age: 30},
				{Name: "Alice", Age: 30},
				{Name: "Bob", Age: 30},
			},
			age:           30,
			expectedCount: 3,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := CountUsersAboveAge(tt.users, tt.age)
			require.Equal(t, tt.expectedCount, result, "Count of users above age does not match expected result")
		})
	}
}
