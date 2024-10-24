package user

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	userRoutes := r.Group("/users")
	{
		userRoutes.GET("/", GetUser)
		userRoutes.POST("/create-user", CreateUser)
		userRoutes.GET("/get-user/:userid", GetUser)
		userRoutes.POST("/login", LogInUser)
		userRoutes.POST("/refresh-token", RefreshToken)
	}
}
