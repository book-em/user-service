package test

import (
	"fmt"
	"testing"

	domain "bookem-user-service/domain"
	service "bookem-user-service/service"

	assert "github.com/stretchr/testify/assert"
)

func TestDelete_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	id := uint(1)
	callerID := uint(1)
	user := defaultUser
	user.ID = id
	mockRepo.On("FindById", id).Return(user, nil)
	mockRepo.On("Delete", id).Return()

	err := svc.Delete(callerID, id)

	assert.NoError(t, err)
}

func TestDelete_TriedToDeleteSomeoneElse(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	id := uint(1)
	callerID := uint(2)
	user := defaultUser
	user.ID = id
	mockRepo.On("FindById", id).Return(user, nil)

	err := svc.Delete(callerID, id)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUnauthorized, err)
}

func TestDelete_UserNotFound(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := service.NewService(mockRepo)

	id := uint(1)
	callerID := uint(1)
	mockRepo.On("FindById", id).Return(nil, fmt.Errorf("user not found"))

	err := svc.Delete(callerID, id)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
}
