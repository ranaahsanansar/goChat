package chat

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ranaahsanansar/gochat/initializers"
	chatDto "github.com/ranaahsanansar/gochat/internal/chat/dto"
	"github.com/ranaahsanansar/gochat/models"
	"github.com/ranaahsanansar/gochat/utils"
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

	var memberDetails models.Users
	result := initializers.DB.First(&memberDetails, "username = ?", userName)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "failed to get user",
		})
		return
	}

	fmt.Println(strconv.FormatUint(uint64(memberDetails.ID), 10))
	result = initializers.DB.Create(&room)

	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(500, gin.H{
			"message": "failed to create room",
		})

	} else {
		fmt.Println("Rooming", strconv.FormatUint(uint64(room.ID), 10))
		var isSuccess = AddMemberToGroupHelper(strconv.FormatUint(uint64(room.ID), 10), strconv.FormatUint(uint64(memberDetails.ID), 10), c)
		if isSuccess {
			c.JSON(200, gin.H{
				"message": "room created",
				"data": gin.H{
					"roomId": room.ID,
				},
			})
		}
	}

}

func AutoCreateGroup(c *gin.Context) {
	var memberId = c.Param("userId")

	// get User name  from Db
	var memberDetails models.Users
	result := initializers.DB.First(&memberDetails, memberId)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "failed to get user",
		})
		return
	}

	var groupName = "Inbox with " + memberDetails.Username

	// Create In DB
	room := models.Groups{
		Name:      groupName,
		CreatedBy: memberId,
	}
	result = initializers.DB.Create(&room)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "failed to create room",
		})
		return
	}

	// add member to the room

	groupMember := models.GroupsMembers{
		GroupID: room.ID,
		UserID:  memberDetails.ID,
	}

	result = initializers.DB.Create(&groupMember)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "failed to create room",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "room created",
	})
}

func AddMemberToGroup(c *gin.Context) {
	var memberId = c.Param("userId")
	var groupId = c.Param("groupId")

	fmt.Println(memberId, groupId)
	var isSuccess = AddMemberToGroupHelper(groupId, memberId, c)

	if isSuccess {
		c.JSON(200, gin.H{
			"message": "Member Added",
		})
	}

}

func GetGroups(c *gin.Context) {
	userName, isExists := c.Get("username")
	if !isExists {
		c.JSON(400, gin.H{
			"message": "Un Protected Request",
		})
		return
	}

	var userDetails models.Users
	result := initializers.DB.First(&userDetails, "username = ?", userName)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "failed to get user",
		})
		return
	}
	userId := userDetails.ID

	fmt.Println("UserId : ", userId)

	var groups []map[string]interface{}
	result = initializers.DB.Raw(utils.GetAllGroupsQuery, userId).Scan(&groups)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "failed to get groups",
			"data": gin.H{
				"groups": groups,
			},
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "groups found",
		"data": gin.H{
			"groups": groups,
		},
	})
}

// add friend by username

// add some one to group
