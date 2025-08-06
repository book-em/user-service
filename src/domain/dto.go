package domain

type UserDTO struct {
	Username string `json:"username" binding:"required,max=30"`
	Password string `json:"password" binding:"required,min=4,max=30"`
	Email    string `json:"email"    binding:"required,max=60,email"`
	Name     string `json:"name"     binding:"required,max=60"`
	Surname  string `json:"surname"  binding:"required,max=60"`
	Address  string `json:"address"  binding:"required,max=150"`
	Role     string `json:"role"     binding:"required,oneof=guest host admin"`
}

type UserUpdateDTO struct {
	Id       int     `json:id"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	Name     *string `json:"name"`
	Surname  *string `json:"surname"`
	Address  *string `json:"address"`
	// Role     *string `json:"role" `
}

type LoginDTO struct {
	UsernameOrEmail string `json:"usernameOrEmail" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

type JWTDTO struct {
	Jwt string `json:"jwt"`
}
