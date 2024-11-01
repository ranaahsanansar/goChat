package userDto

type SearchUserGroupDTO struct {
	KeyWord string `json:"keyword" binding:"required"`
}
