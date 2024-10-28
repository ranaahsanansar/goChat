package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ranaahsanansar/gochat/initializers"
	"github.com/ranaahsanansar/gochat/models"
)

func JoinRoomGuard(c *gin.Context) bool {

	// extract id from param
	var roomid = c.Param("roomId")
	var token = c.Query("token")

	if roomid == "" || token == "" {

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "room id and token are required",
		})
		c.Abort()
		return false
	}

	isActive, err := IntrospectToken(token, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return false
	}

	if !isActive {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return false
	}

	var chatRoom models.Groups
	result := initializers.DB.First(&chatRoom, roomid)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "failed to get room",
		})
		c.Abort()
		return false
	}

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{
			"message": "room not found",
		})
		c.Abort()
		return false
	}

	// TODO: Add more logic to validate user Who can join this room
	return true

}
