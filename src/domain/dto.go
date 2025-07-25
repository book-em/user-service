package domain

type UserDTO struct {
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Email    string  `json:"email" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Surname  string  `json:"surname" binding:"required"`
	Role     string  `json:"role" binding:"required"`
	Address  Address `json:"address" binding:"required"`
}
