package validations

import (
	"github.com/luisfernandomoraes/observability-go/internal/models"
	"fmt"
)

func ValidateUsers(users []models.User) error {
	for i, user := range users {
		if user.Name == "" {
			return fmt.Errorf("user at index %d has an empty name", i)
		}
		if user.Age < 0 || user.Age > 130 {
			return fmt.Errorf("user at index %d has an invalid age: %d", i, user.Age)
		}
	}
	return nil
}
