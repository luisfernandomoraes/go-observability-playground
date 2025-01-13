package services

import (
	"github.com/luisfernandomoraes/observability-go/internal/models"
	"sort"
)

func FilterUsersAboveAge(users []models.User, age int) []models.User {
	filtered := make([]models.User, 0)
	for _, user := range users {
		if user.Age > age {
			filtered = append(filtered, user)
		}
	}
	return filtered
}

func SortUsersByName(users []models.User) []models.User {
	sorted := append([]models.User{}, users...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})
	return sorted
}

func GroupUsersByAge(users []models.User) map[string][]models.User {
	grouped := map[string][]models.User{
		"<20":   {},
		"20-29": {},
		"30-39": {},
		"40+":   {},
	}

	for _, user := range users {
		switch {
		case user.Age < 20:
			grouped["<20"] = append(grouped["<20"], user)
		case user.Age >= 20 && user.Age < 30:
			grouped["20-29"] = append(grouped["20-29"], user)
		case user.Age >= 30 && user.Age < 40:
			grouped["30-39"] = append(grouped["30-39"], user)
		default:
			grouped["40+"] = append(grouped["40+"], user)
		}
	}
	return grouped
}

// UpdateUsersAge increments the age of each user by the given increment.
func UpdateUsersAge(users []models.User, increment int) {
	for i := range users {
		users[i].Age += increment
	}
}

// CountUsersAboveAge counts the number of users above the given age.
func CountUsersAboveAge(users []models.User, age int) int {
	count := 0
	for _, user := range users {
		if user.Age >= age {
			count++
		}
	}
	return count
}
