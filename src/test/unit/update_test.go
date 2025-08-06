package test

import (
	"testing"

	domain "bookem-user-service/domain"
	service "bookem-user-service/service"

	assert "github.com/stretchr/testify/assert"
)

func TestUpdate_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	// Prepare

	newName := "new123"
	oldSurname := "Jones"

	userBefore := domain.User{ID: 1, Username: "user123", Surname: oldSurname}
	dto := domain.UserUpdateDTO{Id: 1, Username: &newName, Surname: nil}

	userAfter := userBefore
	userAfter.Username = newName
	userAfter.Surname = oldSurname // Because dto.Surname == nil

	// Mock

	mockRepo.On("FindById", 1).Return(&userBefore, nil)
	mockRepo.On("Update", &userBefore).Return(nil)

	// Verify

	newUser, err := svc.Update(1, dto)

	assert.NoError(t, err)
	assert.Equal(t, userAfter, *newUser)
}
