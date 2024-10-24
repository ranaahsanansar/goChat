package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ranaahsanansar/gochat/initializers"
	"github.com/ranaahsanansar/gochat/internal/chat"
	"github.com/ranaahsanansar/gochat/internal/user"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	// Module-wise routes
	user.RegisterRoutes(r)
	chat.ChatRoutes(r)

	// Start server
	r.Run(":3000")
}
