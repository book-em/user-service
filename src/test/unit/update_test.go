package test

import (
	"context"
	"fmt"
	"testing"

	domain "bookem-user-service/domain"

	assert "github.com/stretchr/testify/assert"
)

func TestUpdate_Success(t *testing.T) {
	svc, mockRepo, _ := createTestService()

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
	mockRepo.On("FindByUsernameOrEmailNotId", newName, "", dto.Id).Return(nil, nil)
	mockRepo.On("Update", &userBefore).Return(nil)

	// Verify

	newUser, err := svc.Update(context.Background(), 1, dto)

	assert.NoError(t, err)
	assert.Equal(t, userAfter, *newUser)
}

func TestUpdate_SomeoneElse(t *testing.T) {
	svc, _, _ := createTestService()

	// Prepare

	newName := "new123"

	dto := domain.UserUpdateDTO{Id: 123, Username: &newName}

	// Verify

	newUser, err := svc.Update(context.Background(), 1, dto)

	assert.Nil(t, newUser)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUnauthorized, err)
}

func TestUpdate_UserNotFound(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	// Prepare

	newName := "new123"

	dto := domain.UserUpdateDTO{Id: 1, Username: &newName}

	// Mock

	mockRepo.On("FindById", uint(1)).Return(nil, fmt.Errorf("user not found"))

	// Verify

	newUser, err := svc.Update(context.Background(), uint(1), dto)

	assert.Nil(t, newUser)
	assert.Error(t, err)
}

func TestUpdate_UsernameTaken(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	// Prepare

	oldName := "user123"
	newName := "new123"

	userBefore := domain.User{ID: 1, Username: oldName}
	dto := domain.UserUpdateDTO{Id: 1, Username: &newName}

	exitingUser := domain.User{ID: 123, Username: newName, Email: "anything@email.com"}

	// Mock

	mockRepo.On("FindById", uint(1)).Return(&userBefore, nil)
	mockRepo.On("FindByUsernameOrEmailNotId", newName, "", dto.Id).Return(&exitingUser, nil)

	// Verify

	newUser, err := svc.Update(context.Background(), 1, dto)

	assert.Nil(t, newUser)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUsernameExists, err)
}

func TestUpdate_EmailTaken(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	// Prepare

	oldEmail := "user123@email.com"
	newEmail := "new123@email.com"

	userBefore := domain.User{ID: 1, Email: oldEmail}
	dto := domain.UserUpdateDTO{Id: 1, Email: &newEmail}

	exitingUser := domain.User{ID: 123, Email: newEmail, Username: "Anything other than this"}

	// Mock

	mockRepo.On("FindById", uint(1)).Return(&userBefore, nil)
	mockRepo.On("FindByUsernameOrEmailNotId", "", newEmail, dto.Id).Return(&exitingUser, nil)

	// Verify

	newUser, err := svc.Update(context.Background(), 1, dto)

	assert.Nil(t, newUser)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrEmailExists, err)
}

func TestUpdate_UsernameTakenEmailOk(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	// Prepare

	oldName := "user1"
	newName := "user2"
	oldEmail := "user123@email.com"
	newEmail := "user123@email.com"
	okEmail := "completely_different@email.com"

	userBefore := domain.User{ID: 1, Email: oldEmail, Username: oldName}
	dto := domain.UserUpdateDTO{Id: 1, Email: &newEmail, Username: &newName}

	exitingUser := domain.User{ID: 123, Email: okEmail, Username: newName}

	// Mock

	mockRepo.On("FindById", uint(1)).Return(&userBefore, nil)
	mockRepo.On("FindByUsernameOrEmailNotId", newName, newEmail, dto.Id).Return(&exitingUser, nil)

	// Verify

	newUser, err := svc.Update(context.Background(), 1, dto)

	assert.Nil(t, newUser)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUsernameExists, err)
}

func TestUpdate_UsernameOkEmailTaken(t *testing.T) {
	svc, mockRepo, _ := createTestService()

	// Prepare

	oldName := "user1"
	newName := "user1"
	oldEmail := "user123@email.com"
	newEmail := "userNew1234@email.com"
	okName := "completely_different"

	userBefore := domain.User{ID: 1, Email: oldEmail, Username: oldName}
	dto := domain.UserUpdateDTO{Id: 1, Email: &newEmail, Username: &newName}

	exitingUser := domain.User{ID: 123, Email: newEmail, Username: okName}

	// Mock

	mockRepo.On("FindById", uint(1)).Return(&userBefore, nil)
	mockRepo.On("FindByUsernameOrEmailNotId", newName, newEmail, dto.Id).Return(&exitingUser, nil)

	// Verify

	newUser, err := svc.Update(context.Background(), 1, dto)

	assert.Nil(t, newUser)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrEmailExists, err)
}
