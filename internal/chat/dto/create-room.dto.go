package chatDto

type CreateRoomDTO struct {
	Name string `json:"name" binding:"required"`
}
