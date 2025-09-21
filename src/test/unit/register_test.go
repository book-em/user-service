package test

import (
	domain "bookem-user-service/domain"
	"context"
	"errors"
	"strings"
	"testing"

	assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSuccess(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	dto := *defaultUserDTO

	mockRepo.On("FindByUsernameOrEmail", dto.Username, dto.Email).Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := svc.Register(context.Background(), &dto)

	assert.NoError(t, err)
	assert.Equal(t, dto.Username, user.Username)
	assert.Equal(t, strings.ToLower(dto.Email), user.Email)
	assert.Equal(t, domain.UserRole("guest"), user.Role)
	mockRepo.AssertExpectations(t)
}

func TestUsernameExists(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	dto := *defaultUserDTO
	dto.Username = "username"

	existing := *defaultUser
	existing.Username = "username"

	mockRepo.On("FindByUsernameOrEmail", dto.Username, dto.Email).Return(&existing, nil)

	user, err := svc.Register(context.Background(), &dto)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, domain.ErrUsernameExists)
}

func TestEmailExists(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	dto := *defaultUserDTO
	dto.Username = "user1"
	dto.Email = "mail@mail.com"

	existing := *defaultUser
	existing.Username = "user2"
	existing.Email = "mail@mail.com"

	mockRepo.On("FindByUsernameOrEmail", dto.Username, dto.Email).Return(&existing, nil)

	user, err := svc.Register(context.Background(), &dto)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, domain.ErrEmailExists)
}

func TestCreateFails(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	dto := *defaultUserDTO

	mockRepo.On("FindByUsernameOrEmail", dto.Username, dto.Email).Return(nil, nil)
	mockRepo.On("Create", mock.Anything).Return(errors.New("db down"))

	user, err := svc.Register(context.Background(), &dto)

	assert.Nil(t, user)
	assert.ErrorContains(t, err, "db down")
	assert.ErrorIs(t, err, domain.ErrDBInternal)
}
