package domain

type UserDTO struct {
	Username string  `json:"username" binding:"required,max=30"`
	Password string  `json:"password" binding:"required,min=4,max=30"`
	Email    string  `json:"email"    binding:"required,max=60,email"`
	Name     string  `json:"name"     binding:"required,max=60"`
	Surname  string  `json:"surname"  binding:"required,max=60"`
	Role     string  `json:"role"     binding:"required,oneof=guest host admin"`
	Address  Address `json:"address"  binding:"required"`
}
