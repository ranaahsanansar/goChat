package user

import (
	"github.com/ranaahsanansar/gochat/initializers"
	"github.com/ranaahsanansar/gochat/models"
)

func GetUserByUsername(username string) (models.Users, error) {
	var user models.Users
	result := initializers.DB.Where("username = ?", username).First(&user)
	return user, result.Error
}
