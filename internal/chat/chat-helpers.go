package chat

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ranaahsanansar/gochat/initializers"
	"github.com/ranaahsanansar/gochat/models"
)

func AddMemberToGroupHelper(groupId string, memberId string, c *gin.Context) bool {
	// get User name  from Db
	var memberDetails models.Users
	fmt.Println("Member Id : ", memberId)
	result := initializers.DB.First(&memberDetails, memberId)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "failed to get user",
		})
		return false
	}

	parsedGroupID, err := strconv.ParseUint(groupId, 10, 32)
	if err != nil {
		fmt.Println("Error converting groupId to uint:", err)
		return false
	}

	//  check member is not already added
	var checkGroupMember models.GroupsMembers
	result = initializers.DB.First(&checkGroupMember, "group_id = ? AND user_id = ?", parsedGroupID, memberDetails.ID)
	if result.Error == nil {
		c.JSON(500, gin.H{
			"message": "user already added",
		})
		return false
	}

	var groupMember = models.GroupsMembers{
		GroupID: uint(parsedGroupID),
		UserID:  memberDetails.ID,
	}

	result = initializers.DB.Create(&groupMember)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "failed to create room",
		})
		return false
	}

	return true
}
