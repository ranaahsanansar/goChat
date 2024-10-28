package main

import (
	"github.com/ranaahsanansar/gochat/initializers"
	"github.com/ranaahsanansar/gochat/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Users{}, models.Groups{})
}
