package test

import (
	"fmt"
	"testing"

	domain "bookem-user-service/domain"
	service "bookem-user-service/service"
	util "bookem-user-service/util"

	assert "github.com/stretchr/testify/assert"
)

func TestChangePassword_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	// Prepare

	old := "123"
	new := "abc"
	oldHashed, _ := util.HashPassword(old)
	newHashed, _ := util.HashPassword(new)

	userBefore := domain.User{ID: 1, Password: oldHashed}
	dto := domain.PasswordUpdateDTO{Id: 1, OldPassword: old, NewPassword: new, NewPasswordConfirm: new}

	userAfter := userBefore
	userAfter.Password = newHashed

	// Mock

	mockRepo.On("FindById", uint(1)).Return(&userBefore, nil)
	mockRepo.On("Update", &userBefore).Return(nil)

	// Verify

	newUser, err := svc.ChangePassword(1, dto)

	t.Log(oldHashed)

	assert.NoError(t, err)
	assert.True(t, util.VerifyPassword(newUser.Password, new) == nil)
}

func TestChangePassword_SomeoneElse(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	// Prepare

	dto := domain.PasswordUpdateDTO{Id: 123, OldPassword: "", NewPassword: "", NewPasswordConfirm: ""}

	// Verify

	newUser, err := svc.ChangePassword(1, dto)

	assert.Nil(t, newUser)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrUnauthorized, err)
}

func TestChangePassword_UserNotFound(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	// Prepare

	dto := domain.PasswordUpdateDTO{Id: 123, OldPassword: "", NewPassword: "", NewPasswordConfirm: ""}

	// Mock

	mockRepo.On("FindById", uint(1)).Return(nil, fmt.Errorf("user not found"))

	// Verify

	newUser, err := svc.ChangePassword(1, dto)

	assert.Nil(t, newUser)
	assert.Error(t, err)
}

func TestChangePassword_PasswordsNotMatch(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	// Prepare

	old := "123"
	new := "abc"
	oldHashed, _ := util.HashPassword(old)

	userBefore := domain.User{ID: 1, Password: oldHashed}
	dto := domain.PasswordUpdateDTO{Id: 1, OldPassword: old, NewPassword: new, NewPasswordConfirm: ""}

	// Mock

	mockRepo.On("FindById", uint(1)).Return(&userBefore, nil)

	// Verify

	newUser, err := svc.ChangePassword(1, dto)

	assert.Nil(t, newUser)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrPasswordsNotMatch, err)
}

func TestChangePassword_BadOldPassword(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	// Prepare

	old := "123"
	new := "abc"
	oldHashed, _ := util.HashPassword(old)

	userBefore := domain.User{ID: 1, Password: oldHashed}
	dto := domain.PasswordUpdateDTO{Id: 1, OldPassword: "", NewPassword: new, NewPasswordConfirm: new}

	// Mock

	mockRepo.On("FindById", uint(1)).Return(&userBefore, nil)

	// Verify

	newUser, err := svc.ChangePassword(1, dto)

	assert.Nil(t, newUser)
	assert.Error(t, err)
}

func TestChangePassword_PasswordIsTheSame(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	// Prepare

	old := "123"
	oldHashed, _ := util.HashPassword(old)

	userBefore := domain.User{ID: 1, Password: oldHashed}
	dto := domain.PasswordUpdateDTO{Id: 1, OldPassword: old, NewPassword: old, NewPasswordConfirm: old}

	// Mock

	mockRepo.On("FindById", uint(1)).Return(&userBefore, nil)

	// Verify

	newUser, err := svc.ChangePassword(1, dto)

	assert.Nil(t, newUser)
	assert.Error(t, err)
	assert.Equal(t, domain.ErrPasswordNotChanged, err)
}
