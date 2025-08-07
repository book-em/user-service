package test

import (
	domain "bookem-user-service/domain"
	service "bookem-user-service/service"
	"errors"
	"strings"
	"testing"

	assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSuccess(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	dto := *defaultUserDTO

	mockRepo.On("FindByUsernameOrEmail", dto.Username, dto.Email).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := svc.Register(&dto)

	assert.NoError(t, err)
	assert.Equal(t, dto.Username, user.Username)
	assert.Equal(t, strings.ToLower(dto.Email), user.Email)
	assert.Equal(t, domain.UserRole("guest"), user.Role)
	mockRepo.AssertExpectations(t)
}

func TestUsernameExists(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	dto := *defaultUserDTO
	dto.Username = "username"

	existing := *defaultUser
	existing.Username = "username"

	mockRepo.On("FindByUsernameOrEmail", dto.Username, dto.Email).Return(&existing, nil)

	user, err := svc.Register(&dto)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, domain.ErrUsernameExists)
}

func TestEmailExists(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	dto := *defaultUserDTO
	dto.Username = "user1"
	dto.Email = "mail@mail.com"

	existing := *defaultUser
	existing.Username = "user2"
	existing.Email = "mail@mail.com"

	mockRepo.On("FindByUsernameOrEmail", dto.Username, dto.Email).Return(&existing, nil)

	user, err := svc.Register(&dto)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, domain.ErrEmailExists)
}

func TestCreateFails(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	dto := *defaultUserDTO

	mockRepo.On("FindByUsernameOrEmail", dto.Username, dto.Email).Return(nil, nil)
	mockRepo.On("Create", mock.Anything).Return(errors.New("db down"))

	user, err := svc.Register(&dto)

	assert.Nil(t, user)
	assert.ErrorContains(t, err, "db down")
	assert.ErrorIs(t, err, domain.ErrDBInternal)
}
