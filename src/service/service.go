package api

import (
	"bookem-user-service/domain"
	repo "bookem-user-service/repo"
	util "bookem-user-service/util"
	"regexp"
	"strings"
)

type UserService interface {
	Register(input *domain.UserDTO) (*domain.User, error)
}

type userService struct {
	repo repo.UserRepository
}

func NewUserService(r repo.UserRepository) UserService {
	return &userService{r}
}

func (s *userService) Register(dto *domain.UserDTO) (*domain.User, error) {

	if !isValidEmail(dto.Email) {
		return nil, domain.ErrInvalidEmail
	}

	hashed, err := util.HashPassword(dto.Password)
	if err != nil {
		return nil, domain.ErrHashingPassword
	}

	user := &domain.User{
		Username: dto.Username,
		Password: hashed,
		Email:    strings.ToLower(dto.Email),
		Name:     dto.Name,
		Surname:  dto.Surname,
		Role:     domain.UserRole(strings.ToLower(dto.Role)),
		Address:  dto.Address,
	}

	existing := s.repo.FindByUsernameOrEmail(dto.Username, dto.Email)
	if existing != nil {
		if existing.Username == dto.Username {
			return nil, domain.ErrUserExists
		}
		if existing.Email == dto.Email {
			return nil, domain.ErrEmailExists
		}
	}

	return user, nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
