package api

import (
	"bookem-user-service/domain"
	repo "bookem-user-service/repo"
	util "bookem-user-service/util"
	"fmt"
	"strings"
)

type Service interface {
	Register(input *domain.UserDTO) (*domain.User, error)
}

type service struct {
	repo repo.Repository
}

func NewService(r repo.Repository) Service {
	return &service{r}
}

func (s *service) Register(dto *domain.UserDTO) (*domain.User, error) {

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
		Role:     domain.UserRole(dto.Role),
		Address:  dto.Address,
	}

	existing := s.repo.FindByUsernameOrEmail(dto.Username, dto.Email)
	if existing != nil {
		if existing.Username == dto.Username {
			return nil, domain.ErrUsernameExists
		}
		if existing.Email == dto.Email {
			return nil, domain.ErrEmailExists
		}
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrDBInternal, err)
	}

	return user, nil
}
