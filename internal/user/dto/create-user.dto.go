package user

type RegisterUserDTO struct {
	Username        string `json:"username" binding:"required,min=3,max=20"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=2"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=2"`
}
