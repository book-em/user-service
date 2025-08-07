package test

import (
	"fmt"
	"testing"

	domain "bookem-user-service/domain"
	service "bookem-user-service/service"

	assert "github.com/stretchr/testify/assert"
)

func TestUpdate_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	// Prepare

	oldName := "user123"
	newName := "new123"
	oldSurname := "Jones"

	userBefore := domain.User{ID: 1, Username: oldName, Surname: oldSurname}
	dto := domain.UserUpdateDTO{Id: 1, Username: &newName, Surname: nil}

	userAfter := userBefore
	userAfter.Username = newName
	userAfter.Surname = oldSurname // Because dto.Surname == nil

	// Mock

	mockRepo.On("FindById", uint(1)).Return(&userBefore, nil)
	mockRepo.On("FindByUsernameOrEmail", newName, "").Return(nil, nil)
	mockRepo.On("Update", &userBefore).Return(nil)

	// Verify

	newUser, err := svc.Update(1, dto)

	assert.NoError(t, err)
	assert.Equal(t, userAfter, *newUser)
}

func TestUpdate_SomeoneElse(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	// Prepare

	newName := "new123"

	dto := domain.UserUpdateDTO{Id: 123, Username: &newName}

	// Verify

	newUser, err := svc.Update(1, dto)

	assert.Nil(t, newUser)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUnauthorized, err)
}

func TestUpdate_UserNotFound(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	// Prepare

	newName := "new123"

	dto := domain.UserUpdateDTO{Id: 1, Username: &newName}

	// Mock

	mockRepo.On("FindById", uint(1)).Return(nil, fmt.Errorf("user not found"))

	// Verify

	newUser, err := svc.Update(uint(1), dto)

	assert.Nil(t, newUser)
	assert.Error(t, err)
}

func TestUpdate_UsernameTaken(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	// Prepare

	oldName := "user123"
	newName := "new123"

	userBefore := domain.User{ID: 1, Username: oldName}
	dto := domain.UserUpdateDTO{Id: 1, Username: &newName}

	exitingUser := domain.User{ID: 123}

	// Mock

	mockRepo.On("FindById", uint(1)).Return(&userBefore, nil)
	mockRepo.On("FindByUsernameOrEmail", newName, "").Return(&exitingUser, nil)

	// Verify

	newUser, err := svc.Update(1, dto)

	assert.Nil(t, newUser)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUsernameExists, err)
}

func TestUpdate_EmailTaken(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	// Prepare

	oldName := "user123@email.com"
	newName := "new123@email.com"

	userBefore := domain.User{ID: 1, Email: oldName}
	dto := domain.UserUpdateDTO{Id: 1, Email: &newName}

	exitingUser := domain.User{ID: 123}

	// Mock

	mockRepo.On("FindById", uint(1)).Return(&userBefore, nil)
	mockRepo.On("FindByUsernameOrEmail", "", newName).Return(&exitingUser, nil)

	// Verify

	newUser, err := svc.Update(1, dto)

	assert.Nil(t, newUser)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrEmailExists, err)
}
