package chat

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ranaahsanansar/gochat/initializers"
	chatDto "github.com/ranaahsanansar/gochat/internal/chat/dto"
	"github.com/ranaahsanansar/gochat/models"
)

func CreateGroup(c *gin.Context) {
	var body chatDto.CreateRoomDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userName, isExists := c.Get("username")
	if !isExists {
		c.JSON(400, gin.H{
			"message": "Un Protected Request",
		})
		return
	}

	// Create In DB
	room := models.Groups{
		Name:      body.Name,
		CreatedBy: userName.(string),
	}
	initializers.DB.Create(&room)
	c.JSON(200, gin.H{
		"message": "room created",
	})

}
