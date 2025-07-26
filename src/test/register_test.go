package test

import (
	"errors"
	"strings"
	"testing"

	domain "bookem-user-service/domain"
	service "bookem-user-service/service"

	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindByUsernameOrEmail(username, email string) *domain.User {
	args := m.Called(username, email)
	if user, ok := args.Get(0).(*domain.User); ok {
		return user
	}
	return nil
}

func (m *MockUserRepo) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// NOTE: As Gin covers validation, I won’t check for
// nil values, empty values, min/max cardinality, or email type.

var defaultUserDTO = &domain.UserDTO{
	Username: "user",
	Password: "pass",
	Email:    "email@mail.com",
	Name:     "name",
	Surname:  "surname",
	Role:     "guest",
	Address:  "Address 123",
}

var defaultUser = &domain.User{
	Username: "user",
	Password: "pass",
	Email:    "email@mail.com",
	Name:     "name",
	Surname:  "surname",
	Role:     "guest",
	Address:  "Address 123",
}

func TestSuccess(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := service.NewUserService(mockRepo)

	dto := *defaultUserDTO

	mockRepo.On("FindByUsernameOrEmail", dto.Username, dto.Email).Return(nil)
	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := svc.Register(&dto)

	assert.NoError(t, err)
	assert.Equal(t, dto.Username, user.Username)
	assert.Equal(t, strings.ToLower(dto.Email), user.Email)
	assert.Equal(t, domain.UserRole("guest"), user.Role)
	mockRepo.AssertExpectations(t)
}

func TestUsernameExists(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := service.NewUserService(mockRepo)

	dto := *defaultUserDTO
	dto.Username = "username"

	existing := *defaultUser
	existing.Username = "username"

	mockRepo.On("FindByUsernameOrEmail", dto.Username, dto.Email).Return(&existing)

	user, err := svc.Register(&dto)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, domain.ErrUsernameExists)
}

func TestEmailExists(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := service.NewUserService(mockRepo)

	dto := *defaultUserDTO
	dto.Username = "user1"
	dto.Email = "mail@mail.com"

	existing := *defaultUser
	existing.Username = "user2"
	existing.Email = "mail@mail.com"

	mockRepo.On("FindByUsernameOrEmail", dto.Username, dto.Email).Return(&existing)

	user, err := svc.Register(&dto)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, domain.ErrEmailExists)
}

func TestCreateFails(t *testing.T) {
	mockRepo := new(MockUserRepo)
	svc := service.NewUserService(mockRepo)

	dto := *defaultUserDTO

	mockRepo.On("FindByUsernameOrEmail", dto.Username, dto.Email).Return(nil)
	mockRepo.On("Create", mock.Anything).Return(errors.New("db down"))

	user, err := svc.Register(&dto)

	assert.Nil(t, user)
	assert.ErrorContains(t, err, "db down")
	assert.ErrorIs(t, err, domain.ErrDBInternal)
}
