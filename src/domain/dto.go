package domain

type UserCreateDTO struct {
	Username string `json:"username" binding:"required,max=30"`
	Password string `json:"password" binding:"required,min=4,max=30"`
	Email    string `json:"email"    binding:"required,max=60,email"`
	Name     string `json:"name"     binding:"max=60"`
	Surname  string `json:"surname"  binding:"max=60"`
	Address  string `json:"address"  binding:"max=150"`
	Role     string `json:"role"     binding:"required,oneof=guest host admin"`
}

type UserDTO struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"   `
	Name     string `json:"name"    `
	Surname  string `json:"surname" `
	Address  string `json:"address" `
	Role     string `json:"role"    `
}

type UserUpdateDTO struct {
	Id       uint    `json:"id"`
	Username *string `json:"username"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Name     *string `json:"name"`
	Surname  *string `json:"surname"`
	Address  *string `json:"address"`
}

func NewUserDTO(user *User) UserDTO {
	return UserDTO{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
		Surname:  user.Surname,
		Address:  user.Address,
		Role:     string(user.Role),
	}
}

type LoginDTO struct {
	UsernameOrEmail string `json:"usernameOrEmail" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

type PasswordUpdateDTO struct {
	Id                 uint   `json:"id"`
	OldPassword        string `json:"oldPassword"`
	NewPassword        string `json:"newPassword"`
	NewPasswordConfirm string `json:"newPasswordConfirm"`
}

type JWTDTO struct {
	Jwt string `json:"jwt"`
}
