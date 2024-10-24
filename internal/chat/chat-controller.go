package chat

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func JoinRoom(c *gin.Context) {
	// extract id from param
	var roomid = c.Param("roomid")
	fmt.Println(roomid)
	c.JSON(200, gin.H{
		"message": "joined room",
	})
}
