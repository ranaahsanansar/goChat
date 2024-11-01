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
	var memberId = c.Param("username")
	var groupId = c.Param("groupId")

	var userDetails models.Users
	result := initializers.DB.First(&userDetails, "username = ?", memberId)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "failed to get user",
		})
		return
	}
	memberId = strconv.FormatUint(uint64(userDetails.ID), 10)

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
	result = initializers.DB.Raw(utils.GetAllGroupsQuery, userName, userId).Scan(&groups)

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

func DeleteGroup(c *gin.Context) {
	var groupId = c.Param("groupId")
	var userName, isExists = c.Get("username")
	if !isExists {
		c.JSON(400, gin.H{
			"message": "Un Protected Request",
		})
		return
	}

	if groupId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group id is required"})
		c.Abort()
		return
	}

	//  check user is the creator of the group
	var groupDetails models.Groups
	result := initializers.DB.First(&groupDetails, groupId)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "failed to get group",
		})
		return
	}
	if groupDetails.CreatedBy != userName {
		c.JSON(401, gin.H{
			"message": "you are not the creator of this group",
		})
		return
	}

	result = initializers.DB.Unscoped().Delete(&models.Groups{}, groupId)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "failed to delete",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "group deleted",
	})
}

func GetGroupInfo(c *gin.Context) {
	var groupId = c.Param("groupId")

	var groupDetails models.Groups
	result := initializers.DB.First(&groupDetails, groupId)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "failed to get group",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "group found",
		"data": gin.H{
			"group": groupDetails,
		},
	})
}

// add friend by username

// add some one to group
